package services

import (
	"context"
	"fmt"
)

// GetBalance возвращает сумму баллов и сумму списаний.
func (serv *Services) GetBalance(ctx context.Context, userID int) (current, withdraw float64, err error) {
	// Возвращает баланс
	current, err = serv.store.GetCurrentSumAccrual(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("GetCurrentBalance fail: %w", err)
	}
	withdraw, err = serv.store.GetWithDraw(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("GetWithDraw fail: %w", err)
	}
	current -= withdraw
	return current, withdraw, nil
}
