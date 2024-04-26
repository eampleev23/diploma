package services

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/eampleev23/diploma/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func (serv *Services) GetStatusFromAccrual(ctx context.Context,
	textPlainContent string,
	userID int) (
	o models.Order,
	err error) {
	serv.logger.ZL.Debug("GetStatusFromAccrual has started..")
	try := 1
	for o.Status != "PROCESSED" && o.Status != "INVALID" {
		o, err = serv.uploadOrderTry(ctx, textPlainContent, userID)
		if err != nil {
			return models.Order{}, fmt.Errorf("uploadOrderTry fail: %w", err)
		}

		serv.logger.ZL.Debug("Got status from accrual",
			zap.String("status", o.Status),
			zap.String("order", o.Number),
			zap.Int("user ID", o.CustomerID),
			zap.Int("try", try),
		)
		if o.Status == "REGISTERED" {
			o.Status = "NEW"
		}
		orderBack, err := serv.store.UpdateOrder(ctx, o)
		if err != nil {
			return models.Order{}, fmt.Errorf("UpdateOrder fail: %w", err)
		}
		serv.logger.ZL.Debug("UpdateOrder success..",
			zap.String("status", orderBack.Status),
			zap.Int("try", try),
		)
		// Здесь ранее была не корректная задержка на 10 милисекунд, надо пересмотреть насколько она нужна.
		try++
	}
	serv.logger.ZL.Debug("GetStatusFromAccrual has finished..")
	return o, err
}

func (serv *Services) uploadOrderTry(
	ctx context.Context, //nolint:unparam // корректно всегда передавать контекст
	textPlainContent string,
	userID int) (
	o models.Order,
	err error) {
	// Отправляем первый запрос в систему рассчета баллов лояльности
	// создаём новый клиент
	client := resty.New()
	client.
		// устанавливаем количество повторений
		SetRetryCount(3). //nolint:gomnd // no magic
		// длительность ожидания между попытками
		SetRetryWaitTime(30 * time.Second). //nolint:gomnd // no magic
		// длительность максимального ожидания
		SetRetryMaxWaitTime(90 * time.Second) //nolint:gomnd // no magic

	var responseErr models.MyAPIError
	var orderAccrualResp models.OrderAccrualResp

	URL, err := url.JoinPath(serv.config.AccrualRunAddr+"/api/orders/", textPlainContent)
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
	serv.logger.ZL.Debug("got order status from accrual",
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
	return o, nil
}
