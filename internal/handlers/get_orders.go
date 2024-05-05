package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// GetOrders возвращает все заказы пользователя.
func (h *Handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	h.logger.ZL.Debug("GetOrders has started..")
	// Проверяем авторизацию
	// Ппроверяем, не авторизован ли пользователь, отправивший запрос.
	h.logger.ZL.Debug("Checking auth..")
	userID, err := h.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.logger.ZL.Debug("Authorized user:", zap.Int("userID", userID))

	orders, err := h.services.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		h.logger.ZL.Error("GetOrdersByUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		h.logger.ZL.Debug("No data for response")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	ownersOrders, err := models.GetResponseGetOwnerOrders(orders)
	if err != nil {
		h.logger.ZL.Info("GetResponseGetOwnerOrders fail", zap.Error(err))
		return
	}
	if err := enc.Encode(ownersOrders); err != nil {
		h.logger.ZL.Info("fail encoding response in handler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
