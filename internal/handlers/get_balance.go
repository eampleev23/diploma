package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// GetBalance возвращает текущую сумму баллов лояльности и сумму использованных баллов.
func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(keyUserIDCtx).(int)
	if !ok {
		h.logger.ZL.Error("Error getting user ID from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	current, withdraw, err := h.services.GetBalance(r.Context(), userID)
	if err != nil {
		h.logger.ZL.Error("Service GetBalance fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	response := models.GetResponseBalance(current, withdraw)
	if err := enc.Encode(response); err != nil {
		h.logger.ZL.Info("fail encoding response in handler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
