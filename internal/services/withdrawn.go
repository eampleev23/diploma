package services

import (
	"context"
	"fmt"
	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

func (serv *Services) MakeWithdrawn(ctx context.Context, withdrawn models.Withdrawn) (
	success, isOrder, isEnough bool,
	err error) {
	serv.l.ZL.Debug("Service MakeWithdrawn has started..")
	order, err := serv.s.GetFullOrderByOrder(ctx, withdrawn.Order)
	if err != nil {
		return success, isOrder, isEnough, err
	}
	serv.l.ZL.Debug("got GetFullOrderByOrder result",
		zap.Int("order id", order.ID),
		zap.Int("order accrual", order.Accrual),
		zap.String("order status", order.Status),
	)

	isOrder = true

	current, withdrawnSum, err := serv.GetBalance(ctx, withdrawn.UserID)
	if err != nil {
		return success, isOrder, isEnough, fmt.Errorf("serv.GetBalance fail.. %w", err)
	}
	balance := current - withdrawnSum
	serv.l.ZL.Debug("Counted balance",
		zap.Int("balance", balance),
		zap.Int("current", current),
		zap.Int("withdrawnSum", withdrawnSum),
	)
	if balance > withdrawn.Sum {
		isEnough = true
	}
	success, err = serv.s.CreateWithdrawn(ctx, withdrawn)
	return success, isOrder, isEnough, err
}
