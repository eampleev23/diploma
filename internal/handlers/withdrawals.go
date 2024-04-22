package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// Withdrawals отдает информацию о выводе средств.
func (h *Handlers) Withdrawals(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Debug("Withdrawals handler has started..")
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

	withdrawals, err := h.serv.GetWithdrawalsByUserID(r.Context(), userID)
	if err != nil {
		h.l.ZL.Error("GetWithdrawalsByUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		h.l.ZL.Debug("No data for response")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	ownersWithdrawals, err := models.GetResponseGetOwnerWithdrawals(withdrawals)
	if err != nil {
		h.l.ZL.Info("GetResponseGetOwnerWithdrawals fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := enc.Encode(ownersWithdrawals); err != nil {
		h.l.ZL.Info("fail encoding response in handler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
