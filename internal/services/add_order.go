package services

import (
	"context"
	"fmt"

	"github.com/eampleev23/diploma/internal/models"
)

func (serv *Services) AddOrder(ctx context.Context, newOrder models.Order) (
	orderBack models.Order,
	err error) {
	orderBack, err = serv.store.AddNewOrder(ctx, newOrder)
	if err != nil {
		return orderBack, fmt.Errorf("create row fail..%w", err)
	}
	return orderBack, nil
}
