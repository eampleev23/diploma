package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/handlers"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/services"
	"github.com/eampleev23/diploma/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c, err := cnf.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize a new config: %w", err)
	}

	mL, err := mlg.NewZapLogger(c.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to initialize a new logger: %w", err)
	}
	mL.ZL.Debug("Logger success created..")

	au, err := myauth.Initialize(c, mL)
	if err != nil {
		return fmt.Errorf("failed to initialize a new authorizer: %w", err)
	}

	s, err := store.NewStorage(c, mL)
	if err != nil {
		return fmt.Errorf("failed to initialize a new store: %w", err)
	}

	if len(c.DBDSN) != 0 {
		// Отложенно закрываем соединение с бд.
		defer func() {
			if err := s.DBConnClose(); err != nil {
				mL.ZL.Info("store failed to properly close the DB connection")
			}
		}()
	}
	serv := services.NewServices(s, c, mL, *au)
	h, err := handlers.NewHandlers(s, c, mL, *au, *serv)
	if err != nil {
		return fmt.Errorf("handlers constructor's error: %w", err)
	}

	mL.ZL.Info("Running server", zap.String("address", c.RanAddr))
	r := chi.NewRouter()
	r.Use(mL.RequestLogger)
	//r.Use(au.Auth)
	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Authentication)
	r.Post("/api/user/orders", h.UploadOrder)
	r.Get("/api/user/orders", h.GetOrders)
	r.Get("/api/user/balance", h.GetBalance)
	r.Post("/api/user/balance/withdraw", h.Withdrawn)
	r.Get("/api/user/withdrawals", h.Withdrawals)
	err = http.ListenAndServe(c.RanAddr, r)
	if err != nil {
		return fmt.Errorf("ошибка ListenAndServe: %w", err)
	}
	return nil
}
