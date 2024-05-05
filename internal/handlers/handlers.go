package handlers

import (
	"fmt"
	"net/http"

	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/services"
	"github.com/eampleev23/diploma/internal/store"
)

type Handlers struct {
	store      store.Store
	config     *cnf.Config
	logger     *mlg.ZapLog
	authorizer myauth.Authorizer
	services   services.Services
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
		store:      s,
		config:     c,
		logger:     l,
		authorizer: au,
		services:   serv,
	}, nil
}

func (h *Handlers) GetUserID(r *http.Request) (userID int, err error) {
	h.logger.ZL.Debug("GetUserID started.. ")
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, fmt.Errorf("no cookies by name token %w", err)
	}
	userID, err = h.authorizer.GetUserID(cookie.Value)
	if err != nil {
		return 0, fmt.Errorf("h.authorizer.GetUserID fail %w", err)
	}
	return userID, nil
}
