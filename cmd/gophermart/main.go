package main

import (
	"fmt"
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/handlers"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/eampleev23/diploma/internal/myauth"
	"github.com/eampleev23/diploma/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	mL, err := mlg.NewZapLogger("info")
	if err != nil {
		return fmt.Errorf("failed to initialize a new logger: %w", err)
	}
	mL.ZL.Info("Logger success created..")

	c, err := cnf.NewConfig(mL)
	if err != nil {
		return fmt.Errorf("failed to initialize a new config: %w", err)
	}

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

	h, err := handlers.NewHandlers(s, c, mL, *au)
	if err != nil {
		return fmt.Errorf("handlers constructor's error: %w", err)
	}

	mL.ZL.Info("Running server", zap.String("address", c.RanAddr))
	r := chi.NewRouter()
	r.Use(mL.RequestLogger)
	r.Post("/api/user/register", h.Register)
	err = http.ListenAndServe(c.RanAddr, r)
	if err != nil {
		return fmt.Errorf("ошибка ListenAndServe: %w", err)
	}
	return nil
}
