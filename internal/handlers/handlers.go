package handlers

import (
	"fmt"
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/store"
	"net/http"
)

var keyUserIDCtx myauth.Key = myauth.KeyUserIDCtx

type Handlers struct {
	s  store.Store
	c  *cnf.Config
	l  *mlg.ZapLog
	au myauth.Authorizer
}

func NewHandlers(s store.Store, c *cnf.Config, l *mlg.ZapLog, au myauth.Authorizer) (*Handlers, error) {
	return &Handlers{
		s:  s,
		c:  c,
		l:  l,
		au: au,
	}, nil
}

func (h *Handlers) GetUserID(r *http.Request) (userID int, isAuth bool, err error) {
	userIDCtx, ok := r.Context().Value(keyUserIDCtx).(int)
	if !ok {
		return 0, false, fmt.Errorf("userIDCtx is not set: %w", err)
	}
	if userIDCtx != 0 {
		return userIDCtx, false, nil
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, false, fmt.Errorf("token not set in cookie: %w", err)
	}
	userID, err = h.au.GetUserID(cookie.Value)
	if err != nil {
		return 0, false, fmt.Errorf("au.GetUserID error: %w", err)
	}
	return userID, true, nil
}
