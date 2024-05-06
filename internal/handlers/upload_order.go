package handlers

import (
	"context"
	"errors"
	"github.com/eampleev23/diploma/internal/models"
	"github.com/eampleev23/diploma/internal/store"
	"go.uber.org/zap"
	"net/http"
)

// UploadOrder добавляет новый заказ в систему (заявка на получение баллов лояльности).
func (h *Handlers) UploadOrder(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(keyUserIDCtx).(int)
	if !ok {
		h.logger.ZL.Error("Fail getting userID from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	textPlainContent, err := h.services.GetTextPlain(r)
	if err != nil {
		h.logger.ZL.Error("GetTextPlain fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.ZL.Debug("", zap.String("textPlainContent", textPlainContent))
	// Далее проверяем алгоритмом луна и возвращаем 422 если не верный формат номера заказа.
	err = h.services.MoonCheck(textPlainContent)
	if err != nil {
		h.logger.ZL.Debug("MoonTest fail", zap.Error(err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	h.logger.ZL.Debug("Moon test success..")

	newOrder := models.Order{
		Number:     textPlainContent,
		CustomerID: userID,
		Status:     "NEW",
	}

	_, err = h.services.AddOrder(r.Context(), newOrder)
	if err != nil && errors.Is(err, store.ErrConflict) {
		h.logger.ZL.Debug("this order already exists")
		confUserID, _ := h.services.GetUserIDByOrderNumber(r.Context(), textPlainContent)
		if confUserID == userID {
			// Заказ был создан текущим пользователем
			w.WriteHeader(http.StatusOK)
			return
		}
		// Заказ был загружен другим пользователем
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		h.logger.ZL.Debug("AddOrder fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go h.getFromAccrual(r.Context(), textPlainContent, userID)
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) getFromAccrual(ctx context.Context, textPlainContent string, userID int) {
	_, err := h.services.GetStatusFromAccrual(ctx, textPlainContent, userID)
	if err != nil {
		h.logger.ZL.Debug("getFromAccrual fail..", zap.Error(err))
	}
}
