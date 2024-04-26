package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// Withdrawn списывает баллы лояльности.
func (h *Handlers) Withdrawn(w http.ResponseWriter, r *http.Request) {
	h.logger.ZL.Debug("Withdrawn handler has started..")
	// Проверяем формат запроса
	contentType := r.Header.Get("Content-Type")
	supportsJSON := strings.Contains(contentType, "application/json")
	if !supportsJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Проверяем авторизацию
	// Ппроверяем, не авторизован ли пользователь, отправивший запрос.
	h.logger.ZL.Debug("Checking auth..")
	userID, isAuth, err := h.GetUserID(r)
	if err != nil {
		h.logger.ZL.Error("GetUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAuth {
		h.logger.ZL.Debug("Unauthorized user..")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.logger.ZL.Debug("Authorized user:", zap.Int("userID", userID))
	var req models.Withdrawn
	// Декодер работает потоково, кажется это правильнее + короче, чем анмаршал.
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		h.logger.ZL.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.UserID = userID
	h.logger.ZL.Debug("got request",
		zap.String("order", req.Order),
		zap.Float64("sum", req.Sum),
	)
	err = h.services.MoonCheck(req.Order)
	if err != nil {
		h.logger.ZL.Debug("Mooncheck fail..")
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	if err := h.services.MakeWithdrawn(r.Context(), req); err != nil {
		h.logger.ZL.Info("Sum of the balance is not enough..", zap.Error(err))
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	h.logger.ZL.Debug("Success debit")
}
