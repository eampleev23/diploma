package handlers

import (
	"encoding/json"
	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Withdrawn списывает баллы лояльности
func (h *Handlers) Withdrawn(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Debug("Withdrawn handler has started..")
	// Проверяем формат запроса
	contentType := r.Header.Get("Content-Type")
	supportsJSON := strings.Contains(contentType, "application/json")
	if !supportsJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Проверяем авторизацию
	// Ппроверяем, не авторизован ли пользователь, отправивший запрос.
	h.l.ZL.Debug("Checking auth..")
	userID, isAuth, err := h.GetUserID(r)
	if err != nil {
		h.l.ZL.Error("GetUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAuth {
		h.l.ZL.Debug("Unauthorized user..")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.l.ZL.Debug("Authorized user:", zap.Int("userID", userID))
	var req models.Withdrawn
	// Декодер работает потоково, кажется это правильнее + короче, чем анмаршал.
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		h.l.ZL.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.l.ZL.Debug("got request",
		zap.String("order", req.Order),
		zap.Int("sum", req.Sum),
	)
}
