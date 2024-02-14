package store

import (
	"context"
	"fmt"
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/models"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type Store interface {
	// DBConnClose закрывает соединение с базой данных
	DBConnClose() (err error)
	// InsertUser добавляет нового пользователя или добавляет ошибку о конфликте данных
	InsertUser(ctx context.Context, userRegReq models.UserRegReq) (userBack models.User, err error)
}

func NewStorage(c *cnf.Config, l *mlg.ZapLog) (Store, error) {
	s, err := NewDBStore(c, l)
	if err != nil {
		return nil, fmt.Errorf("error creating new db store: %w", err)
	}
	l.ZL.Info("DB store created success..")
	return s, nil
}
