package services

import (
	"context"
	"fmt"
)

// GetBalance возвращает сумму баллов и сумму списаний
func (serv *Services) GetBalance(ctx context.Context, userID int) (current, withdraw float64, err error) {
	// Возвращает баланс
	serv.l.ZL.Debug("services / GetBalance started..")
	current, err = serv.s.GetCurrentSumAccrual(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("GetCurrentBalance fail: %w", err)
	}
	withdraw, err = serv.s.GetWithDraw(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("GetWithDraw fail: %w", err)
	}
	current = current - withdraw
	return current, withdraw, nil
}
