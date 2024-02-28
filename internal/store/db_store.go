package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
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

func (D DBStore) DBConnClose() (err error) {
	if err := D.dbConn.Close(); err != nil {
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
func (D DBStore) InsertUser(ctx context.Context, userRegReq models.UserRegReq) (newUser models.User, err error) {
	newUser = models.User{}
	err = D.dbConn.QueryRow("INSERT INTO users (login, password) VALUES($1, $2) RETURNING id, login, password", userRegReq.Login, userRegReq.Password).Scan(&newUser.ID, &newUser.Login, &newUser.Password)
	// Проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = ErrConflict
	}
	return newUser, err
}
func (D DBStore) GetUserByLoginAndPassword(ctx context.Context, userLoginReq models.UserLoginReq) (userBack models.User, err error) {
	userBack = models.User{}
	row := D.dbConn.QueryRowContext(ctx,
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

func (D DBStore) AddNewOrder(ctx context.Context, newOrder models.Order) (orderBack models.Order, err error) {
	orderBack = models.Order{}
	err = D.dbConn.QueryRow("INSERT INTO orders (number, customer_id) VALUES($1, $2) RETURNING id, number, customer_id",
		newOrder.Number,
		newOrder.UserID).Scan(
		&orderBack.ID,
		&orderBack.Number,
		&orderBack.UserID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = ErrConflict
	}
	return orderBack, err
}

func (D DBStore) GetUserIDByOrder(ctx context.Context, orderNumber string) (userID int, err error) {
	row := D.dbConn.QueryRowContext(ctx,
		`SELECT customer_id FROM orders WHERE number = $1 LIMIT 1`,
		orderNumber,
	)
	err = row.Scan(&userID) // Разбираем результат
	if err != nil {
		return userID, fmt.Errorf("faild to get user id by order's number %w", err)
	}
	return userID, nil
}
