package myauth

import (
	"context"
	"fmt"
	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
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
		_, err := r.Cookie("token")
		if err != nil {
			// Получаем логгер из контекста запроса
			logger, ok := r.Context().Value(keyLogger).(*mlg.ZapLog)
			if !ok {
				log.Printf("Error getting logger")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			logger.ZL.Debug("No cookie", zap.String("err", err.Error()))
			// Cookie не установлена, устанавливаем
			//newRandomUserID, err := au.setNewCookie(w)
			//if err != nil {
			//	logger.ZL.Info("Error setting cookie:", zap.String("err", err.Error()))
			//}
			//logger.ZL.Debug("Success setted cookie", zap.Int("newRandomUserId", newRandomUserID))
			//ctx := context.WithValue(r.Context(), KeyUserIDCtx, newRandomUserID)
			//logger.ZL.Debug("Передали newRandomUserID", zap.Int("newRandomUserID", newRandomUserID))
			//next.ServeHTTP(w, r.WithContext(ctx))
			next.ServeHTTP(w, r)
			return
		}
		// если кука уже установлена, то через контекст передаем 0
		ctx := context.WithValue(r.Context(), KeyUserIDCtx, 0)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (au *Authorizer) SetNewCookie(w http.ResponseWriter, userID int) (err error) {

	au.l.ZL.Debug("setNewCookie got userID", zap.Int("userID", userID))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(au.c.TokenExp)),
		},
		// Собственное утверждение
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(au.c.SecretKey))
	if err != nil {
		return fmt.Errorf("token.SignedString fail.. %w", err)
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	http.SetCookie(w, &cookie)
	return nil
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
