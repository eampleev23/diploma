package services

import (
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/store"
)

type Services struct {
	s  store.Store
	c  *cnf.Config
	l  *mlg.ZapLog
	au myauth.Authorizer
}

func NewServices(s store.Store, c *cnf.Config, l *mlg.ZapLog, au myauth.Authorizer) Services {
	services := Services{
		s:  s,
		c:  c,
		l:  l,
		au: au,
	}
	return services
}
