package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// GetOrders возвращает все заказы пользователя.
func (h *Handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(keyUserIDCtx).(int)
	if !ok {
		h.logger.ZL.Error("Error getting user ID from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	orders, err := h.services.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		h.logger.ZL.Error("GetOrdersByUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
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
