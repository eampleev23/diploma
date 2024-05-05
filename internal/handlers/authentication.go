package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// Authentication авторизует зарегистрированного пользователя.
func (h *Handlers) Authentication(w http.ResponseWriter, r *http.Request) {
	// Проверяем формат запроса
	contentType := r.Header.Get("Content-Type")                       //nolint:goconst //not needed
	supportsJSON := strings.Contains(contentType, "application/json") //nolint:goconst //not needed
	if !supportsJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Получаем данные в случае корректного запроса.
	var req models.UserLoginReq
	// Декодер работает потоково, кажется это правильнее + короче, чем анмаршал.
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		h.logger.ZL.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Проверяем, не авторизован ли пользователь, отправивший запрос.
	userIDAlreadyAuth, err := h.GetUserID(r)
	if err != nil {
		h.logger.ZL.Error("GetUserID fail") //nolint:goconst //not needed
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.ZL.Debug("isAuth", zap.Bool("auth", isAuth))

	authUser, err := h.store.GetUserByLoginAndPassword(r.Context(), req)
	if err != nil {
		h.logger.ZL.Info("User is not found", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if isAuth && userIDAlreadyAuth != authUser.ID {
		h.logger.ZL.Error("Already authorized user trying to login by another one", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.authorizer.SetNewCookie(w, authUser.ID)
	if err != nil {
		h.logger.ZL.Error("SetNewCookie fail", zap.Error(err))
		w.WriteHeader(http.StatusOK)
		return
	}
	h.logger.ZL.Debug("Success authorization, user id -", zap.Int("authUser.ID", authUser.ID))
	h.logger.ZL.Debug("Success authorization, user login -", zap.String("authUser.Login", authUser.Login))
}
