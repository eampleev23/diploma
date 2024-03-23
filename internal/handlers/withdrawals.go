package handlers

import "net/http"

// Withdrawals отдает информацию о выводе средств
func (h *Handlers) Withdrawals(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Debug("Withdrawals handler has started..")
}
