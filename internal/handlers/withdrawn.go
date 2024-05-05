package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eampleev23/diploma/internal/models"
	"go.uber.org/zap"
)

// Withdrawn списывает баллы лояльности.
func (h *Handlers) Withdrawn(w http.ResponseWriter, r *http.Request) {
	// Проверяем формат запроса
	contentType := r.Header.Get("Content-Type")
	supportsJSON := strings.Contains(contentType, "application/json")
	if !supportsJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(keyUserIDCtx).(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var req models.Withdrawn
	// Декодер работает потоково, кажется это правильнее + короче, чем анмаршал.
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		h.logger.ZL.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.UserID = userID
	err := h.services.MoonCheck(req.Order)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	if err := h.services.MakeWithdrawn(r.Context(), req); err != nil {
		h.logger.ZL.Info("Sum of the balance is not enough..", zap.Error(err))
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
}
