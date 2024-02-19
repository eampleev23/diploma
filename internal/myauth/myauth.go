package myauth

import (
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"log"
	"net/http"
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

// Claims описывает утверждения, хранящиеся в токене + добавляет кастомное UserID.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// GetUserID возвращает ID пользователя.
func (au *Authorizer) GetUserID(tokenString string) (int, error) {
	// Создаем экземпляр структуры с утверждениями
	claims := &Claims{}
	// Парсим из строки токена tokenString в структуру claims
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(au.c.SecretKey), nil
	})
	if err != nil {
		au.l.ZL.Info("Failed in case to get ownerId from token ", zap.Error(err))
	}

	return claims.UserID, nil
}
