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
	appConfig, err := cnf.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize a new config: %w", err)
	}

	logger, err := mlg.NewZapLogger(appConfig.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to initialize a new logger: %w", err)
	}
	logger.ZL.Debug("Logger success created..")

	authorizer, err := myauth.Initialize(appConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize a new authorizer: %w", err)
	}

	appStorage, err := store.NewStorage(appConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize a new store: %w", err)
	}

	if len(appConfig.DBDSN) != 0 {
		// Отложенно закрываем соединение с бд.
		defer func() {
			if err := appStorage.DBConnClose(); err != nil {
				logger.ZL.Info("store failed to properly close the DB connection")
			}
		}()
	}
	appServices := services.NewServices(appStorage, appConfig, logger, authorizer)
	h, err := handlers.NewHandlers(appStorage, appConfig, logger, authorizer, appServices)
	if err != nil {
		return fmt.Errorf("handlers constructor's error: %w", err)
	}

	logger.ZL.Info("Running server", zap.String("address", appConfig.RanAddr))
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Authentication)
	r.Post("/api/user/orders", h.UploadOrder)
	r.Get("/api/user/orders", h.GetOrders)
	r.Get("/api/user/balance", h.GetBalance)
	r.Post("/api/user/balance/withdraw", h.Withdrawn)
	r.Get("/api/user/withdrawals", h.Withdrawals)
	err = http.ListenAndServe(appConfig.RanAddr, r)
	if err != nil {
		return fmt.Errorf("ошибка ListenAndServe: %w", err)
	}
	return nil
}
