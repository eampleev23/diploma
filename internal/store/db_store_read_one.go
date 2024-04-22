package store

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
)

func (d DBStore) GetUserByLoginAndPassword(
	ctx context.Context,
	userLoginReq models.UserLoginReq,
) (
	userBack models.User,
	err error) {
	userBack = models.User{}
	row := d.dbConn.QueryRowContext(ctx,
		`SELECT id, login, password FROM users WHERE login = $1 AND password = $2 LIMIT 1`,
		userLoginReq.Login,
		userLoginReq.Password,
	)
	err = row.Scan(&userBack.ID, &userBack.Login, &userBack.Password) // Разбираем результат
	if err != nil {
		return userBack, fmt.Errorf("faild to get user by login and password like this %w", err)
	}
	return userBack, nil
}
func (d DBStore) GetUserIDByOrder(ctx context.Context, orderNumber string) (userID int, err error) {
	row := d.dbConn.QueryRowContext(ctx,
		`SELECT customer_id FROM orders WHERE number = $1 LIMIT 1`,
		orderNumber,
	)
	err = row.Scan(&userID) // Разбираем результат
	if err != nil {
		return userID, fmt.Errorf("faild to get user id by order's number %w", err)
	}
	return userID, nil
}

func (d DBStore) GetFullOrderByOrder(
	ctx context.Context,
	orderNumber string) (
	fullOrder models.Order,
	err error) {
	d.l.ZL.Debug("GetFullOrderByOrder has started..")
	row := d.dbConn.QueryRowContext(ctx,
		`SELECT 
    			id,number,customer_id,status,accrual,uploaded_at
				FROM orders
				WHERE number = $1 LIMIT 1`,
		orderNumber,
	)
	err = row.Scan(&fullOrder.ID, &fullOrder.Number,
		&fullOrder.CustomerID, &fullOrder.Status,
		&fullOrder.Accrual, &fullOrder.UploadedAt) // Разбираем результат
	if err != nil {
		return fullOrder, fmt.Errorf("faild to get full order by order's number %w", err)
	}
	return fullOrder, err
}
