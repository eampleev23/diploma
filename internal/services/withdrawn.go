package services

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

func (serv *Services) MakeWithdrawn(ctx context.Context, withdrawn models.Withdrawn) (err error) {
	serv.l.ZL.Debug("Service MakeWithdrawn2 has started..")
	serv.l.ZL.Debug("Got withdrawn",
		zap.Float64("sum", withdrawn.Sum),
		zap.Int("user_id", withdrawn.UserID),
	)
	// Здесь получаем баланс в виде суммы начислений и суммы списаний
	current, withdrawnSum, err := serv.GetBalance(ctx, withdrawn.UserID)
	if err != nil {
		return fmt.Errorf("serv.GetBalance fail.. %w", err)
	}
	serv.l.ZL.Debug("Получаем баланс в виде суммы начислений и суммы списаний",
		zap.Float64("current", current),
		zap.Float64("withdrawnSum", withdrawnSum),
	)
	serv.l.ZL.Debug("Баланс правильно посчитался")
	serv.l.ZL.Debug("Дальше нужно понять хватает ли баланса для снятия")
	if current-withdrawnSum < 0 {
		// На балансе недостаточно денег
		return fmt.Errorf("недостаточно баллов для списания")
	} else {
		withdrawnBack, err := serv.s.CreateWithdrawn(ctx, withdrawn)
		if err != nil {
			return fmt.Errorf("CreateWithdrawn fail: %w", err)
		}
		serv.l.ZL.Debug("WithdrawnBack",
			zap.Int("ID", withdrawnBack.ID),
			zap.Float64("SUM", withdrawnBack.Sum),
			zap.Int("USER", withdrawnBack.UserID),
			zap.String("ORDER", withdrawnBack.Order),
		)
		return nil
	}
}
