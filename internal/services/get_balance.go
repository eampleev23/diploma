package services

import (
	"context"
)

// GetBalance возвращает сумму баллов и сумму списаний
func (serv *Services) GetBalance(ctx context.Context, userID int) (current, withdraw int, err error) {
	// Возвращает заказы пользователя по ID
	serv.l.ZL.Debug("services / GetBalance started..")
	return current, withdraw, nil
}
