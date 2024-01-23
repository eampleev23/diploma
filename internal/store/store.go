package store

import (
	"fmt"

	"github.com/eampleev23/diploma.git/cmd/internal/cnf"
	"github.com/eampleev23/diploma.git/cmd/internal/mlg"
)

type Store interface {
	// Close закрывает соединение с базой данных
	Close() (err error)
}

func NewStorage(c *cnf.Config, l *mlg.ZapLog) (Store, error) {
	s, err := NewDBStore(c, l)
	if err != nil {
		return nil, fmt.Errorf("error creating new db store: %w", err)
	}
	l.ZL.Info("DB store created success..")
	return s, nil
}
