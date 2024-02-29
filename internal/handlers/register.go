package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// Register регистрирует нового пользователя.
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	// Проверяем формат запроса
	contentType := r.Header.Get("Content-Type")
	supportsJSON := strings.Contains(contentType, "application/json")
	if !supportsJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем, не авторизован ли пользователь, отправивший запрос.
	_, isAuth, err := h.GetUserID(r)
	if err != nil {
		h.l.ZL.Error("GetUserID fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.l.ZL.Debug("isAuth", zap.Bool("auth", isAuth))
	if isAuth {
		h.l.ZL.Error("already authorized user is going to register")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получаем данные в случае корректного запроса.
	var req models.UserRegReq
	// Декодер работает потоково, кажется это правильнее + короче, чем анмаршал.
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		h.l.ZL.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser, err := h.s.InsertUser(r.Context(), req)
	if err != nil {
		h.l.ZL.Debug("login is not uniq", zap.Error(err))
		w.WriteHeader(http.StatusConflict)
		return
	}
	// мы здесь если пользователь успешно зарегистрирован
	// надо установить куку
	// а в самом начале надо проверить на куку, возможно он уже авторизован и тогда надо отправлять
	// внуреннюю ошибку сервера

	err = h.au.SetNewCookie(w, newUser.ID)
	if err != nil {
		h.l.ZL.Debug("Not authorize after success register", zap.Error(err))
		w.WriteHeader(http.StatusOK)
		return
	}
}
