package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Register регистрирует нового пользователя
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	supportsJSON := strings.Contains(contentType, "application/json")
	if !supportsJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// проверяем, не авторизован ли пользователь, отправивший запрос
	userID, _, err := h.GetUserID(r)
	// ...
	fmt.Println(userID)

	// Получаем данные в случае корректного запроса.
	var req models.UserRegReq
	// Декодер работает потоково, кажется это правильнее + короче, чем анмаршал.
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		h.l.ZL.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userBack, err := h.s.InsertUser(r.Context(), req)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	// мы здесь если пользователь успешно зарегистрирован
	// надо установить куку
	// а в самом начале надо проверить на куку, возможно он уже авторизован и тогда надо отправлять
	// внуреннюю ошибку сервера

	fmt.Println("userBack=", userBack)
}
