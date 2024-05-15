package store

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
)

func (d DBStore) UpdateOrder(ctx context.Context, o models.Order) (orderBack models.Order, err error) {
	orderBack = models.Order{}
	err = d.dbConn.QueryRow( //nolint:execinquery // нужен скан
		`UPDATE orders SET status = $1, accrual = $2 WHERE number = $3
				RETURNING
    			id, status, accrual;`,
		o.Status,
		o.Accrual,
		o.Number,
	).Scan(
		&orderBack.ID,
		&orderBack.Status,
		&orderBack.Accrual)
	if err != nil {
		return models.Order{}, fmt.Errorf("d.dbConn.QueryRow fail %w", err)
	}
	return orderBack, nil
}
