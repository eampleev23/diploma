package store

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

func (d DBStore) GetCurrentSumAccrual(ctx context.Context, userID int) (current float64, err error) {
	d.l.ZL.Debug("DBStore / GetCurrentSumAccrual has started..")
	row := d.dbConn.QueryRowContext(ctx,
		`SELECT SUM(accrual)
				FROM orders
				WHERE customer_id = $1;`, userID)
	err = row.Scan(&current) // Разбираем результат
	if err != nil {
		return current, fmt.Errorf("faild to get sum accrual by scan %w", err)
	}
	d.l.ZL.Debug("DBStore / GetCurrentSumAccrual / Got sum accrual", zap.Float64("current", current))
	return current, nil
}

func (d DBStore) GetWithDraw(ctx context.Context, userID int) (withdraw float64, err error) {
	d.l.ZL.Debug("DBStore / GetWithDraw has started..")
	row := d.dbConn.QueryRowContext(ctx,
		`SELECT SUM(sum)
				FROM withdraw
				WHERE user_id = $1;`, userID)
	err = row.Scan(&withdraw) // Разбираем результат
	if err != nil {
		strError := error.Error(err)
		if strings.Contains(strError, "converting NULL to int is unsupported") {
			return 0, nil
		}
		return 0, fmt.Errorf("QueryRowContext fail: %w", err)
	}
	return withdraw, nil
}
