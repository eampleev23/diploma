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

	current, withdrawnSum, err := serv.GetBalance(ctx, withdrawn.UserID)
	if err != nil {
		return success, isOrder, isEnough, fmt.Errorf("serv.GetBalance fail.. %w", err)
	}
	balance := current - withdrawnSum
	serv.l.ZL.Debug("Counted balance",
		zap.Float64("balance", balance),
		zap.Float64("current", current),
		zap.Float64("withdrawnSum", withdrawnSum),
	)
	if balance >= withdrawn.Sum {
		isEnough = true
	} else {
		isEnough = false
		return success, isOrder, isEnough, err
	}
	success, withdrawnBack, err := serv.s.CreateWithdrawn(ctx, withdrawn)
	if err != nil {
		return success, isOrder, isEnough, fmt.Errorf("CreateWithdrawn fail: %w", err)
	}
	serv.l.ZL.Debug("WithdrawnBack",
		zap.Int("ID", withdrawnBack.ID),
		zap.Float64("SUM", withdrawnBack.Sum),
		zap.Int("USER", withdrawnBack.UserID),
		zap.String("ORDER", withdrawnBack.Order),
	)
	return success, isOrder, isEnough, err
}
