package services

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
)

func (serv *Services) GetOrdersByUserID(
	ctx context.Context, userID int) (
	orders []models.Order, err error) {
	// Возвращает заказы пользователя по ID
	orders, err = serv.store.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return orders, fmt.Errorf("store method fail: %w", err)
	}
	return orders, nil
}
