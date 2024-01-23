package handlers

import (
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/store"
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
