package store

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

func (d DBStore) GetOrdersByUserID(ctx context.Context, userID int) (orders []models.Order, err error) {
	d.l.ZL.Debug("db_store / GetOrdersByUserID started..")
	rows, err := d.dbConn.QueryContext( //nolint:sqlclosecheck // not clear
		ctx,
		`SELECT 
    				id,number,customer_id,status,accrual,uploaded_at
					FROM
					    orders
					WHERE
					customer_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("not get orders for customer by customer_id %w", err)
	}
	for rows.Next() {
		var v models.Order
		err = rows.Scan(&v.ID, &v.Number, &v.CustomerID, &v.Status, &v.Accrual, &v.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf(" rows san fail: %w", err)
		}
		d.l.ZL.Debug("got order",
			zap.String("number", v.Number),
			zap.Time("uploaded at", v.UploadedAt),
			zap.Int("customer", v.CustomerID),
			zap.String("status", v.Status),
			zap.Float64("accrual", v.Accrual),
			zap.Int("id", v.ID),
		)
		orders = append(orders, v)
	}
	// проверяем на ошибки
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err in db store GetOrdersByUserID: %w", err)
	}
	return orders, nil
}

func (d DBStore) GetWithdrawalsByUserID(
	ctx context.Context,
	userID int) (
	withdeawals []models.Withdrawn,
	err error) {
	d.l.ZL.Debug("db_store / GetWithdrawalsByUserID started..")
	rows, err := d.dbConn.QueryContext( //nolint:sqlclosecheck // not clear
		ctx,
		`SELECT 
    				id, sum, order_number, user_id, processed_at
					FROM
					    withdraw
					WHERE
					user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("not get orders for customer by customer_id %w", err)
	}
	for rows.Next() {
		var v models.Withdrawn
		err = rows.Scan(&v.ID, &v.Sum, &v.Order, &v.UserID, &v.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf(" rows scan fail: %w", err)
		}
		d.l.ZL.Debug("got withdrawn",
			zap.Int("id", v.ID),
			zap.Float64("sum", v.Sum),
			zap.String("order", v.Order),
			zap.Int("user_id", v.UserID),
			zap.Time("processed at", v.ProcessedAt),
		)
		withdeawals = append(withdeawals, v)
	}
	// проверяем на ошибки
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err in db store GetOrdersByUserID: %w", err)
	}
	return withdeawals, nil
}
