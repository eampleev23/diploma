package handlers

import (
	"net/http"
)

// Register регистрирует нового пользователя
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Info("Registration handler got request..")
}
