package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/eampleev23/diploma/internal/store"

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
	w.WriteHeader(http.StatusAccepted)
}
