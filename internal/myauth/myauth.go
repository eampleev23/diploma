package myauth

import (
	"log"
	"net/http"

	"github.com/eampleev23/diploma.git/cmd/internal/cnf"
	"github.com/eampleev23/diploma.git/cmd/internal/mlg"
)

type Authorizer struct {
	l *mlg.ZapLog
	c *cnf.Config
}

var keyLogger mlg.Key = mlg.KeyLoggerCtx

// Initialize инициализирует синглтон авторизовывальщика с секретным ключом.
func Initialize(c *cnf.Config, l *mlg.ZapLog) (*Authorizer, error) {
	au := &Authorizer{
		c: c,
		l: l,
	}
	return au, nil
}

type Key string

const (
	KeyUserIDCtx Key = "user_id_ctx"
)

// Auth мидлвар, который проверяет авторизацию
func (au *Authorizer) Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Получаем логгер из контекста запроса
		l, ok := r.Context().Value(keyLogger).(*mlg.ZapLog)
		if !ok {
			log.Printf("Error getting logger")
			return
		}
		l.ZL.Info("Empty auth middleware worked..")

	}
	return http.HandlerFunc(fn)
}
