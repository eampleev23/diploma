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
	fmt.Println("userBack=", userBack)
}
