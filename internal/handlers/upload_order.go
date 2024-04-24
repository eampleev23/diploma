package handlers

import (
	"context"
	"fmt"
	"github.com/eampleev23/diploma/internal/models"
	"github.com/eampleev23/diploma/internal/store"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// UploadOrder добавляет новый заказ в систему (заявка на получение баллов лояльности).
func (h *Handlers) UploadOrder(w http.ResponseWriter, r *http.Request) {
	if err := h.uploadOrder(r); err != nil {
		h.l.ZL.Error("failed to handle request to uploadOrder", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) uploadOrder(r *http.Request) error {
	// Проверяем формат запроса.
	contentType := r.Header.Get("Content-Type")
	textPlain := strings.Contains(contentType, "text/plain")
	if !textPlain {
		return fmt.Errorf("format doesn't content text/plain")
	}

	// Ппроверяем, не авторизован ли пользователь, отправивший запрос.
	userID, isAuth, err := h.GetUserID(r)
	if err != nil {
		return fmt.Errorf("h.GetUserID fail.. %w", err)
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

	newOrder := models.Order{
		Number:     textPlainContent,
		CustomerID: userID,
		Status:     "NEW",
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
	if err != nil {
		h.l.ZL.Debug("AddOrder fail..", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go h.getFromAccrual(r.Context(), textPlainContent, userID)
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) getFromAccrual(ctx context.Context, textPlainContent string, userID int) {
	_, err := h.serv.GetStatusFromAccrual(ctx, textPlainContent, userID)
	if err != nil {
		h.l.ZL.Debug("getFromAccrual fail..", zap.Error(err))
	}
}
