package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/handlers"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/services"
	"github.com/eampleev23/diploma/internal/store"
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
	appHandlers, err := handlers.NewHandlers(appStorage, appConfig, logger, authorizer, appServices)
	if err != nil {
		return fmt.Errorf("handlers constructor's error: %w", err)
	}
	logger.ZL.Info("Running server", zap.String("address", appConfig.RanAddr))
	routes := appHandlers.GetRoutes()
	err = http.ListenAndServe(appConfig.RanAddr, routes)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ошибка ListenAndServe: %w", err)
	}
	return nil
}
