package myauth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"go.uber.org/zap"
)

type Authorizer struct {
	l *mlg.ZapLog
	c *cnf.Config
}

// Initialize инициализирует синглтон авторизовывальщика с секретным ключом.
func Initialize(c *cnf.Config, l *mlg.ZapLog) (Authorizer, error) {
	au := Authorizer{
		c: c,
		l: l,
	}
	return au, nil
}

type Key string

const (
	KeyUserIDCtx Key = "user_id_ctx"
)

// Auth мидлвар, который проверяет авторизацию.
func (au *Authorizer) Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Здесь в перспективе будет мидлвар проверки авторизации, надо разобраться
		// как использовать только в определенных хэндлерах
	}
	return http.HandlerFunc(fn)
}

func (au *Authorizer) SetNewCookie(w http.ResponseWriter, userID int) (err error) {
	au.l.ZL.Debug("setNewCookie got userID", zap.Int("userID", userID))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Когда создан токен.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(au.c.TokenExp)),
		},
		// Собственное утверждение.
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
	// Создаем экземпляр структуры с утверждениями.
	claims := &Claims{}
	// Парсим из строки токена tokenString в структуру claims.
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(au.c.SecretKey), nil
	})
	if err != nil {
		au.l.ZL.Info("Failed in case to get ownerId from token ", zap.Error(err))
	}

	return claims.UserID, nil
}
