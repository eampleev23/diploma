package handlers

import "net/http"

// GetBalance возвращает текущую сумму баллов лояльности и сумму использованных баллов
func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	h.l.ZL.Debug("Handler GetBalance has started..")
}
