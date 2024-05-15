package services

import (
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/store"
)

type Services struct {
	store      store.Store
	config     *cnf.Config
	logger     *mlg.ZapLog
	authorizer myauth.Authorizer
}

func NewServices(s store.Store, c *cnf.Config, l *mlg.ZapLog, au myauth.Authorizer) Services {
	services := Services{
		store:      s,
		config:     c,
		logger:     l,
		authorizer: au,
	}
	return services
}
