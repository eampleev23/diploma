package handlers

import (
	"net/http"

	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/services"
	"github.com/eampleev23/diploma/internal/store"
)

type Handlers struct {
	s    store.Store
	c    *cnf.Config
	l    *mlg.ZapLog
	au   myauth.Authorizer
	serv services.Services
}

func NewHandlers(
	s store.Store,
	c *cnf.Config,
	l *mlg.ZapLog,
	au myauth.Authorizer,
	serv services.Services) (
	*Handlers,
	error) {
	return &Handlers{
		s:    s,
		c:    c,
		l:    l,
		au:   au,
		serv: serv,
	}, nil
}

func (h *Handlers) GetUserID(r *http.Request) (userID int, isAuth bool, err error) {
	h.l.ZL.Debug("GetUserID started.. ")
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, false, err
	}
	userID, err = h.au.GetUserID(cookie.Value)
	if err != nil {
		return 0, false, nil
	}
	return userID, true, nil
}
