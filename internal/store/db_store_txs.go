package store

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/eampleev23/diploma/internal/models"
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
	result, err := tx.Exec(`INSERT INTO withdraw
				(sum, order_number, user_id)
				VALUES($1, $2, $3)`, withdrawn.Sum, withdrawn.Order, withdrawn.UserID)
	if err != nil {
		return fmt.Errorf("tx.Exec fail: %w", err)
	}
	lastID, err := result.LastInsertId()
	d.l.ZL.Debug("Fail getting last id for insert MakeWithdrawTX exec", zap.Int64("withdraw id", lastID))

	accrualsSumRaw := tx.QueryRowContext(ctx, `SELECT SUM(accrual)
				FROM orders
				WHERE customer_id = $1;`, withdrawn.UserID)
	if err != nil {
		return fmt.Errorf("accrualsSumRaw queryRowContext fail: %w", err)
	}
	err = accrualsSumRaw.Scan(&accrualsSum) // Получили сумму начислений
	if err != nil {
		return fmt.Errorf("accrualsSumRaw Scan fail: %w", err)
	}
	withdrawalsSumRaw := tx.QueryRowContext(ctx, `SELECT SUM(sum)
				FROM withdraw
				WHERE user_id = $1;`, withdrawn.UserID)
	if err != nil {
		return fmt.Errorf("withdrawalsSumRaw tx.QueryRowContext fail: %w", err)
	}
	err = withdrawalsSumRaw.Scan(&withdrawalsSum) // Получили сумму списаний
	if err != nil {
		log.Println("error b scan")
	}
	if balance := accrualsSum - withdrawalsSum; balance < 0 {
		return fmt.Errorf("forbidden")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx commit for create withdrawals fail: %w", err)
	}
	return nil
}
