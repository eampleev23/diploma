package store

import (
	"context"
	"github.com/eampleev23/diploma/internal/models"
)

func (d DBStore) UpdateOrder(ctx context.Context, o models.Order) (orderBack models.Order, err error) {
	orderBack = models.Order{}
	err = d.dbConn.QueryRow( //nolint:execinquery // нужен скан
		`UPDATE orders SET status = $1 WHERE number = $2
				RETURNING
    			id, status;`,
		o.Status,
		o.Number,
	).Scan(
		&orderBack.ID,
		&orderBack.Status)
	return orderBack, err
}
