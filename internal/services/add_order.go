package services

import (
	"context"
	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

func (serv *Services) AddOrder(ctx context.Context, newOrder models.Order) (err error) {
	serv.l.ZL.Debug("AddOrder start..")
	serv.l.ZL.Debug("AddOrder / got new order number:", zap.String("number", newOrder.Number))
	serv.l.ZL.Debug("AddOrder / got customer_id:", zap.Int("userID", newOrder.UserId))
	orderBack, _ := serv.s.AddNewOrder(ctx, newOrder)
	serv.l.ZL.Debug("AddOrder / got order back id:", zap.Int("orderBackID", orderBack.ID))
	return nil
}
