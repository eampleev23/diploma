package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// Authentication авторизует зарегистрированного пользователя
func (h *Handlers) Authentication(w http.ResponseWriter, r *http.Request) {
	// Проверяем формат запроса
	contentType := r.Header.Get("Content-Type")
	supportsJSON := strings.Contains(contentType, "application/json")
	if !supportsJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Получаем данные в случае корректного запроса.
	var req models.UserLoginReq
	// Декодер работает потоково, кажется это правильнее + короче, чем анмаршал.
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		h.l.ZL.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//проверяем, не авторизован ли пользователь, отправивший запрос
	userIDAlreadyAuth, isAuth, err := h.GetUserID(r)
	if err != nil {
		h.l.ZL.Error("GetUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.l.ZL.Debug("isAuth", zap.Bool("auth", isAuth))

	authUser, err := h.s.GetUserByLoginAndPassword(r.Context(), req)
	if err != nil {
		h.l.ZL.Info("User is not found", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if isAuth && userIDAlreadyAuth != authUser.ID {
		h.l.ZL.Error("Already authorized user trying to login by another one", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.au.SetNewCookie(w, authUser.ID)
	if err != nil {
		h.l.ZL.Error("SetNewCookie fail", zap.Error(err))
		w.WriteHeader(http.StatusOK)
		return
	}
	h.l.ZL.Debug("Success authorization, user id -", zap.Int("authUser.ID", authUser.ID))
	h.l.ZL.Debug("Success authorization, user login -", zap.String("authUser.Login", authUser.Login))
}
