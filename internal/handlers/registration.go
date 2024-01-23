package handlers

import (
	"net/http"
)

// Registration регистрирует нового пользователя
func (h *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Info("Registration handler got request..")
}
