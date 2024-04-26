package services

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

func (serv *Services) AddOrder(ctx context.Context, newOrder models.Order) (
	orderBack models.Order,
	err error) {
	serv.logger.ZL.Debug("AddOrder start..")
	serv.logger.ZL.Debug("AddOrder / got new order number:", zap.String("number", newOrder.Number))
	serv.logger.ZL.Debug("AddOrder / got customer_id:", zap.Int("customerID", newOrder.CustomerID))
	orderBack, err = serv.store.AddNewOrder(ctx, newOrder)
	if err != nil {
		return orderBack, fmt.Errorf("create row fail..%w", err)
	}
	serv.logger.ZL.Debug("AddOrder / got order back id:", zap.Int("orderBackID", orderBack.ID))
	return orderBack, nil
}
