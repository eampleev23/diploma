package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrConflict ошибка, которую используем для сигнала о нарушении целостности данных.
var ErrConflict = errors.New("data conflict")

// InsertUser занимается непосредственно запросом вставки строки в бд.
func (d DBStore) InsertUser(ctx context.Context, userRegReq models.UserRegReq) (newUser models.User, err error) {
	newUser = models.User{}
	err = d.dbConn.QueryRow( //nolint:execinquery // нужен скан
		`INSERT INTO
    users (login, password)
	VALUES($1, $2)
	RETURNING
	    id, login, password`,
		userRegReq.Login,
		userRegReq.Password).Scan(
		&newUser.ID,
		&newUser.Login,
		&newUser.Password)
	// Проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = ErrConflict
	}
	return newUser, err
}

func (d DBStore) AddNewOrder(ctx context.Context, newOrder models.Order) (orderBack models.Order, err error) {
	orderBack = models.Order{}
	err = d.dbConn.QueryRow( //nolint:execinquery // нужен скан
		`INSERT INTO orders
    			(number, customer_id, status, accrual)
				VALUES($1, $2, $3, $4)
				RETURNING
    			id, number, customer_id`,
		newOrder.Number,
		newOrder.CustomerID,
		newOrder.Status,
		newOrder.Accrual).Scan(
		&orderBack.ID,
		&orderBack.Number,
		&orderBack.CustomerID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = ErrConflict
	}
	return orderBack, err
}

func (d DBStore) CreateWithdrawn(
	ctx context.Context,
	withdrawn models.Withdrawn) (
	withdrawnBack models.Withdrawn,
	err error) {
	err = d.dbConn.QueryRow( //nolint:execinquery // нужен скан
		`INSERT INTO withdraw
				(sum, order_number, user_id)
				VALUES($1, $2, $3)
				RETURNING
				id,sum,order_number,user_id`,
		withdrawn.Sum,
		withdrawn.Order,
		withdrawn.UserID).Scan(
		&withdrawnBack.ID,
		&withdrawnBack.Sum,
		&withdrawnBack.Order,
		&withdrawnBack.UserID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = ErrConflict
	}
	if err != nil {
		return withdrawnBack, fmt.Errorf("QueryRow fail: %w", err)
	}
	return withdrawnBack, err
}
