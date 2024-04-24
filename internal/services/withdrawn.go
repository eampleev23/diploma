package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

func (serv *Services) MakeWithdrawn(ctx context.Context, withdrawn models.Withdrawn) (
	success, isOrder, isEnough bool,
	err error) {
	serv.l.ZL.Debug("Service MakeWithdrawn has started..")

	// Здесь получаем баланс в виде суммы начислений и суммы списаний
	current, withdrawnSum, err := serv.GetBalance(ctx, withdrawn.UserID)
	if err != nil {
		return success, isOrder, isEnough, fmt.Errorf("serv.GetBalance fail.. %w", err)
	}
	// Здесь теряются копейки, поэтому переводим в decimal
	cDec := decimal.NewFromFloat(current)
	wDec := decimal.NewFromFloat(withdrawnSum)
	balanceDec := cDec.Sub(wDec)
	balanceDecStr := balanceDec.String()
	balance, err := strconv.ParseFloat(balanceDecStr, 64)
	if err != nil {
		return false, false, false, fmt.Errorf("ParseFloat fail: %w", err)
	}
	fmt.Println("balanceDec=", balanceDec)
	fmt.Println("balance=", balance)
	serv.l.ZL.Debug("Counted balance",
		zap.Float64("balance", balance),
		zap.Float64("current", current),
		zap.Float64("withdrawnSum", withdrawnSum),
	)
	if balance >= withdrawn.Sum {
		isEnough = true
	} else {
		isEnough = false
		return success, isOrder, isEnough, nil
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
	return success, isOrder, isEnough, nil
}
