package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// UploadOrder добавляет новый заказ в систему (заявка на получение баллов лояльности)
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
}
