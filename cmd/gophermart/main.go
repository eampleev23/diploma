package main

import (
	"fmt"
	"github.com/eampleev23/diploma.git/cmd/internal/logger"
	"log"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	myLog, err := logger.NewZapLogger("info")
	if err != nil {
		return fmt.Errorf("failed to initialize a new logger: %w", err)
	}
	myLog.ZL.Info("Logger success created..")
	//
	//c, err := config.NewConfig(myLog)
	//if err != nil {
	//	return fmt.Errorf("failed to initialize a new config: %w", err)
	//}
	//
	//au, err := myauth.Initialize(c.SecretKey, c.TokenEXP, myLog)
	//if err != nil {
	//	return fmt.Errorf("failed to initialize a new authorizer: %w", err)
	//}
	//
	//s, err := store.NewStorage(c, myLog)
	//if err != nil {
	//	return fmt.Errorf("failed to initialize a new store: %w", err)
	//}
	//
	//if len(c.DBDSN) != 0 {
	//	// Отложенно закрываем соединение с бд.
	//	defer func() {
	//		if err := s.Close(); err != nil {
	//			myLog.ZL.Info("store failed to properly close the DB connection")
	//		}
	//	}()
	//}
	//
	//h := handlers.NewHandlers(s, c, myLog, *au)
	//
	//myLog.ZL.Info("Running server", zap.String("address", c.RanAddr))
	//r := chi.NewRouter()
	//err = http.ListenAndServe(c.RanAddr, r)
	//if err != nil {
	//	return fmt.Errorf("ошибка ListenAndServe: %w", err)
	//}
	return nil
}
