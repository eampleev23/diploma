package handlers

import "net/http"

// Withdrawn списывает баллы лояльности
func (h *Handlers) Withdrawn(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Debug("Withdrawn handler has started..")
}
