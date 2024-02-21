package handlers

import (
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
	h.l.ZL.Debug("GetUserID started.. ")
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, false, nil
	}
	userID, err = h.au.GetUserID(cookie.Value)
	if err != nil {
		return 0, false, nil
	}
	return userID, true, nil
}
