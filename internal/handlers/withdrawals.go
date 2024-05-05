package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// Withdrawals отдает информацию о выводе средств.
func (h *Handlers) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(keyUserIDCtx).(int)
	if !ok {
		h.logger.ZL.Error("Fail getting user ID from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	withdrawals, err := h.services.GetWithdrawalsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.ZL.Error("GetWithdrawalsByUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	ownersWithdrawals, err := models.GetResponseGetOwnerWithdrawals(withdrawals)
	if err != nil {
		h.logger.ZL.Info("GetResponseGetOwnerWithdrawals fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := enc.Encode(ownersWithdrawals); err != nil {
		h.logger.ZL.Info("fail encoding response in handler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
