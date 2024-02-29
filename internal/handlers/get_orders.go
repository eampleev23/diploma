package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

// GetOrders возвращает все заказы пользователя.
func (h *Handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Debug("GetOrders has started..")
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

}
