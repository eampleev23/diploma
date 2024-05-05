package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

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

var keyUserIDCtx myauth.Key = myauth.KeyUserIDCtx

func (h *Handlers) GetRoutes() (routes *chi.Mux) {
	routes = chi.NewRouter()

	routes.Group(func(r chi.Router) {
		r.Use(h.logger.RequestLogger)
		r.Group(func(r chi.Router) {
			r.Route("/api/user", func(r chi.Router) {
				r.Use(h.authorizer.Auth)
				r.Method("POST", "/orders", http.HandlerFunc(h.UploadOrder))
				r.Method("GET", "/orders", http.HandlerFunc(h.GetOrders))
				r.Method("GET", "/balance", http.HandlerFunc(h.GetBalance))
				r.Method("POST", "/balance/withdraw", http.HandlerFunc(h.Withdrawn))
				r.Method("GET", "/withdrawals", http.HandlerFunc(h.Withdrawals))
			})
		})
		r.Group(func(r chi.Router) {
			r.Route("/api/user/register", func(r chi.Router) {
				r.Method("POST", "/", http.HandlerFunc(h.Register))
			})
		})
		r.Group(func(r chi.Router) {
			r.Route("/api/user/login", func(r chi.Router) {
				r.Method("POST", "/", http.HandlerFunc(h.Authentication))
			})
		})

	})
	return routes
}

func (h *Handlers) GetUserID(r *http.Request) (userID int, err error) {
	userID, err = h.authorizer.GetUserID(r)
	if err != nil {
		return 0, fmt.Errorf("h.authorizer.GetUserID fail %w", err)
	}
	return userID, nil
}
