package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// GetBalance возвращает текущую сумму баллов лояльности и сумму использованных баллов.
func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Debug("Handler GetBalance has started..")
	// Проверяем авторизацию
	// Ппроверяем, не авторизован ли пользователь, отправивший запрос.
	h.l.ZL.Debug("Checking auth..") //nolint:goconst //not needed
	userID, isAuth, err := h.GetUserID(r)
	if err != nil {
		h.l.ZL.Error("GetUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAuth {
		h.l.ZL.Debug("Unauthorized user..") //nolint:goconst //not needed
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.l.ZL.Debug("Authorized user:", zap.Int("userID", userID)) //nolint:goconst //not needed

	current, withdraw, err := h.serv.GetBalance(r.Context(), userID)
	if err != nil {
		h.l.ZL.Error("Service GetBalance fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.l.ZL.Debug("got balance",
		zap.Float64("current", current),
		zap.Float64("withdraw", withdraw),
	)
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	response, err := models.GetResponseBalance(current, withdraw)
	if err != nil {
		h.l.ZL.Info("GetResponseGetOwnerOrders fail", zap.Error(err))
		return
	}
	if err := enc.Encode(response); err != nil {
		h.l.ZL.Info("fail encoding response in handler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
