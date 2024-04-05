package handlers

import (
	"errors"
	"github.com/eampleev23/diploma/internal/store"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// UploadOrder добавляет новый заказ в систему (заявка на получение баллов лояльности).
func (h *Handlers) UploadOrder(w http.ResponseWriter, r *http.Request) {
	// Проверяем формат запроса.
	contentType := r.Header.Get("Content-Type")
	textPlain := strings.Contains(contentType, "text/plain")
	if !textPlain {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Ппроверяем, не авторизован ли пользователь, отправивший запрос.
	userID, isAuth, err := h.GetUserID(r)
	if err != nil {
		h.l.ZL.Error("GetUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAuth {
		h.l.ZL.Debug("Unaithorized user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.l.ZL.Debug("isAuth", zap.Bool("auth", isAuth))
	h.l.ZL.Debug("", zap.Int("userID", userID))
	textPlainContent, err := h.serv.GetTextPlain(r)
	if err != nil {
		h.l.ZL.Error("GetTextPlain fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.l.ZL.Debug("", zap.String("textPlainContent", textPlainContent))
	// Далее проверяем алгоритмом луна и возвращаем 422 если не верный формат номера заказа.
	err = h.serv.MoonCheck(textPlainContent)
	if err != nil {
		h.l.ZL.Debug("MoonTest fail", zap.Error(err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	h.l.ZL.Debug("Moon test success..")

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
	URL, err := url.JoinPath(h.c.AccrualRunAddr, "/api/good")
	if err != nil {
		h.l.ZL.Debug("url.JoinPath fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"match":"Bork","reward":10,"reward_type":"%"}`).
		Post(URL)
	if err != nil {
		h.l.ZL.Debug("1 req to accrual fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	URL, err = url.JoinPath(h.c.AccrualRunAddr, "/api/orders")
	if err != nil {
		h.l.ZL.Debug("url.JoinPath fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(models.OrderAccrual{
			Order: textPlainContent,
			Goods: []models.Good{{Description: "Чайник Bork", Price: 7000}},
		}).Post(URL)
	if err != nil {
		h.l.ZL.Debug("2 req to accrual fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var responseErr models.MyApiError
	var orderAccrualResp models.OrderAccrualResp

	URL, err = url.JoinPath(h.c.AccrualRunAddr+"/api/orders/", textPlainContent)
	if err != nil {
		h.l.ZL.Debug("url.JoinPath fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = client.R().
		SetError(&responseErr).
		SetResult(&orderAccrualResp).
		Get(URL)

	if err != nil {
		h.l.ZL.Debug("3 req (get) to accrual fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.l.ZL.Debug("got order status",
		zap.String("order", orderAccrualResp.Order),
		zap.String("status", orderAccrualResp.Status),
		zap.Int("accrual", orderAccrualResp.Accrual),
	)

	newOrder := models.Order{
		Number:     orderAccrualResp.Order,
		CustomerID: userID,
		Status:     orderAccrualResp.Status,
		Accrual:    orderAccrualResp.Accrual,
	}
	_, err = h.serv.AddOrder(r.Context(), newOrder)
	if err != nil && errors.Is(err, store.ErrConflict) {
		h.l.ZL.Debug("this order already exists")
		confUserID, _ := h.serv.GetUserIDByOrderNumber(r.Context(), textPlainContent)
		if confUserID == userID {
			// Заказ был создан текущим пользователем
			w.WriteHeader(http.StatusOK)
			return
		}
		// Заказ был загружен другим пользователем
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
