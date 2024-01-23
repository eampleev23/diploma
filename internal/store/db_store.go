package store

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/eampleev23/diploma.git/cmd/internal/cnf"
	"github.com/eampleev23/diploma.git/cmd/internal/mlg"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
)

type DBStore struct {
	dbConn *sql.DB
	c      *cnf.Config
	l      *mlg.ZapLog
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

func (ds DBStore) Close() error {
	if err := ds.dbConn.Close(); err != nil {
		return fmt.Errorf("failed to properly close the DB connection %w", err)
	}
	return nil
}
