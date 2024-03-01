package services

import (
	"context"
	"github.com/eampleev23/diploma/internal/models"
)

func (serv *Services) GetOrdersByUserID(
	ctx context.Context, userID int) (
	orders []models.Order, err error) {
	// Возвращает заказы пользователя по ID
	serv.l.ZL.Debug("services / GetOrdersByUserID started..")
	orders, err = serv.s.GetOrdersByUserID(ctx, userID)
	return orders, nil
}
