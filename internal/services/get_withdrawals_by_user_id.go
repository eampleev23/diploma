package services

import (
	"context"
	"fmt"
	"github.com/eampleev23/diploma/internal/models"
)

func (serv *Services) GetWithdrawalsByUserID(
	ctx context.Context, userID int) (
	withdrawals []models.Withdrawn, err error) {
	// Возвращает заказы пользователя по ID
	serv.l.ZL.Debug("services / GetWithdrawalsByUserID started..")
	withdrawals, err = serv.s.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return withdrawals, fmt.Errorf("store method fail: %w", err)
	}
	return withdrawals, nil
}