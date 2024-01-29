package handlers

import (
	"net/http"
	"strings"
)

// Register регистрирует нового пользователя
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	// Сначала получим контент тайп и отдадим соответствующий ответ если это не JSON
	contentType := r.Header.Get("Content-Type")
	supportsJson := strings.Contains(contentType, "application/json")
	if !supportsJson {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
