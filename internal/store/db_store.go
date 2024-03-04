package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStore struct {
	dbConn *sql.DB
	c      *cnf.Config
	l      *mlg.ZapLog
}

func (d DBStore) DBConnClose() (err error) {
	if err := d.dbConn.Close(); err != nil {
		return fmt.Errorf("failed to properly close the DB connection %w", err)
	}
	return nil
}

func NewDBStore(c *cnf.Config, l *mlg.ZapLog) (*DBStore, error) {
	db, err := sql.Open("pgx", c.DBDSN)
	if err != nil {
		return &DBStore{}, fmt.Errorf("%w", errors.New("sql.open failed in case to create store"))
	}
	if err := runMigrations(c.DBDSN); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}
	return &DBStore{
		dbConn: db,
		c:      c,
		l:      l,
	}, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

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

func (d DBStore) AddNewOrder(ctx context.Context, newOrder models.Order) (orderBack models.Order, err error) {
	orderBack = models.Order{}
	err = d.dbConn.QueryRow( //nolint:execinquery // нужен скан
		`INSERT INTO orders
    			(number, customer_id, status)
				VALUES($1, $2, $3)
				RETURNING
    			id, number, customer_id`,
		newOrder.Number,
		newOrder.CustomerID,
		newOrder.Status).Scan(
		&orderBack.ID,
		&orderBack.Number,
		&orderBack.CustomerID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = ErrConflict
	}
	return orderBack, err
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

func (d DBStore) GetOrdersByUserID(ctx context.Context, userID int) (orders []models.Order, err error) {
	d.l.ZL.Debug("db_store / GetOrdersByUserID started..")
	rows, err := d.dbConn.QueryContext(
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
			return nil, fmt.Errorf(" rows san fail: %w")
		}
		d.l.ZL.Debug("got order",
			zap.String("number", v.Number),
			zap.Time("uploaded at", v.UploadedAt),
			zap.Int("customer", v.CustomerID),
			zap.String("status", v.Status),
			zap.Int("accrual", v.Accrual),
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
