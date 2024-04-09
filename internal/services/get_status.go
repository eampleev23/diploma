package services

import (
	"context"
	"fmt"
	"github.com/eampleev23/diploma/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/url"
	"time"
)

func (serv *Services) GetStatusFromAccrual(ctx context.Context, textPlainContent string, userID int) (o models.Order, err error) {
	serv.l.ZL.Debug("GetStatusFromAccrual has started..")
	for o.Status != "PROCESSED" {
		o, err = serv.uploadOrderTry(ctx, textPlainContent, userID)
		if err != nil {
			return models.Order{}, fmt.Errorf("uploadOrderTry fail: %w", err)
		}
		time.Tick(10)
	}
	serv.l.ZL.Debug("GetStatusFromAccrual has finished..")
	return o, err
}

func (serv *Services) uploadOrderTry(ctx context.Context, textPlainContent string, userID int) (o models.Order, err error) {

	// Отправляем первый запрос в систему рассчета баллов лояльности
	// создаём новый клиент
	client := resty.New()
	client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(30 * time.Second).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(90 * time.Second)

	URL, err := url.JoinPath(serv.c.AccrualRunAddr, "/api/good")
	if err != nil {
		return models.Order{}, fmt.Errorf("url.JoinPath fail: %w", err)
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"match":"Bork","reward":10,"reward_type":"%"}`).
		Post(URL)

	if err != nil {
		return models.Order{}, fmt.Errorf("/api/good request fail: %w", err)
	}
	URL, err = url.JoinPath(serv.c.AccrualRunAddr, "/api/orders")
	if err != nil {
		return models.Order{}, fmt.Errorf("url.JoinPath fail: %w", err)
	}
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(models.OrderAccrual{
			Order: textPlainContent,
			Goods: []models.Good{{Description: "Лейка Bork", Price: 5000}},
		}).Post(URL)
	if err != nil {
		return models.Order{}, fmt.Errorf("second request fail: %w", err)
	}

	var responseErr models.MyAPIError
	var orderAccrualResp models.OrderAccrualResp

	URL, err = url.JoinPath(serv.c.AccrualRunAddr+"/api/orders/", textPlainContent)
	if err != nil {
		return models.Order{}, fmt.Errorf("url.JoinPath fail: %w", err)
	}

	_, err = client.R().
		SetError(&responseErr).
		SetResult(&orderAccrualResp).
		Get(URL)

	if err != nil {
		return models.Order{}, fmt.Errorf("therd request fail: %w", err)
	}
	serv.l.ZL.Debug("got order status",
		zap.String("order", orderAccrualResp.Order),
		zap.String("status", orderAccrualResp.Status),
		zap.Float64("accrual", orderAccrualResp.Accrual),
	)
	o = models.Order{
		Number:     orderAccrualResp.Order,
		CustomerID: userID,
		Status:     orderAccrualResp.Status,
		Accrual:    orderAccrualResp.Accrual,
	}

	return o, err
}
