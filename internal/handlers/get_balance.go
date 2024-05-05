package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// GetBalance возвращает текущую сумму баллов лояльности и сумму использованных баллов.
func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	h.logger.ZL.Debug("Handler GetBalance has started..")
	// Проверяем авторизацию
	// Ппроверяем, не авторизован ли пользователь, отправивший запрос.
	h.logger.ZL.Debug("Checking auth..") //nolint:goconst //not needed
	userID, err := h.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.logger.ZL.Debug("Authorized user:", zap.Int("userID", userID)) //nolint:goconst //not needed

	current, withdraw, err := h.services.GetBalance(r.Context(), userID)
	if err != nil {
		h.logger.ZL.Error("Service GetBalance fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.ZL.Debug("got balance",
		zap.Float64("current", current),
		zap.Float64("withdraw", withdraw),
	)
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	response, err := models.GetResponseBalance(current, withdraw)
	if err != nil {
		h.logger.ZL.Info("GetResponseGetOwnerOrders fail", zap.Error(err))
		return
	}
	if err := enc.Encode(response); err != nil {
		h.logger.ZL.Info("fail encoding response in handler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
