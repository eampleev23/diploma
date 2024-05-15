package services

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

func (serv *Services) MakeWithdrawn(ctx context.Context, withdrawn models.Withdrawn) (err error) {
	// Здесь получаем баланс в виде суммы начислений и суммы списаний
	current, withdrawnSum, err := serv.GetBalance(ctx, withdrawn.UserID)
	if err != nil {
		return fmt.Errorf("serv.GetBalance fail.. %w", err)
	}
	serv.logger.ZL.Debug("Получаем баланс в виде суммы начислений и суммы списаний",
		zap.Float64("current", current),
		zap.Float64("withdrawnSum", withdrawnSum),
	)

	if current-withdrawnSum < 0 {
		// На балансе недостаточно денег
		return fmt.Errorf("недостаточно баллов для списания")
	} else {
		withdrawnBack, err := serv.store.CreateWithdrawn(ctx, withdrawn)
		if err != nil {
			return fmt.Errorf("CreateWithdrawn fail: %w", err)
		}
		serv.logger.ZL.Debug("WithdrawnBack",
			zap.Int("ID", withdrawnBack.ID),
			zap.Float64("SUM", withdrawnBack.Sum),
			zap.Int("USER", withdrawnBack.UserID),
			zap.String("ORDER", withdrawnBack.Order),
		)
		return nil
	}
}

func (serv *Services) MakeWithdrawn1(ctx context.Context, withdrawn models.Withdrawn) (err error) {
	if err = serv.store.MakeWithdrawTX(ctx, withdrawn); err != nil {
		return fmt.Errorf("MakeWithdrawTX fail: %w", err)
	}
	return nil
}
