package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

// Withdrawals отдает информацию о выводе средств
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
		h.l.ZL.Error("GetOrdersByUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		h.l.ZL.Debug("No data for response")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	//enc := json.NewEncoder(w)
	//w.Header().Set("content-type", "application/json")
	//ownersWithdrawals, err := models.GetResponseGetOwnerWithdrawals(withdrawals)
}
