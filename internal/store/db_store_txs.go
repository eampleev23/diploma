package store

import (
	"context"
	"fmt"
	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
	"log"
)

func (d DBStore) MakeWithdrawTX(ctx context.Context, withdrawn models.Withdrawn) (err error) {
	var accrualsSum float64
	var withdrawalsSum float64
	tx, err := d.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("dbConn.BeginTx fail: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()
	tx.Exec(`INSERT INTO withdraw
				(sum, order_number, user_id)
				VALUES($1, $2, $3)`, withdrawn.Sum, withdrawn.Order, withdrawn.UserID)

	a := tx.QueryRowContext(ctx, `SELECT SUM(accrual)
				FROM orders
				WHERE customer_id = $1;`, withdrawn.UserID)
	if err != nil {
		log.Println("error tx.Query at the a")
	}
	err = a.Scan(&accrualsSum) // Получили сумму начислений
	if err != nil {
		d.l.ZL.Error("error a.Scan", zap.Error(err))
	}
	log.Println("accrualsSum=", accrualsSum)
	b := tx.QueryRowContext(ctx, `SELECT SUM(sum)
				FROM withdraw
				WHERE user_id = $1;`, withdrawn.UserID)
	if err != nil {
		log.Println("error b tx query")
	}
	err = b.Scan(&withdrawalsSum) // Получили сумму списаний
	if err != nil {
		log.Println("error b scan")
	}
	log.Println("withdrawalsSum=", withdrawalsSum)
	log.Println("accrualsSum - withdrawalsSum=", accrualsSum-withdrawalsSum)
	if balance := accrualsSum - withdrawalsSum; balance < 0 {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	err = tx.Commit()
	if err != nil {
		log.Println("error commit")
	}
	return err
}
