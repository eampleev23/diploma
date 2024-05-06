package myauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/golang-jwt/jwt/v4"

	"github.com/eampleev23/diploma/internal/cnf"
	"github.com/eampleev23/diploma/internal/mlg"
	"go.uber.org/zap"
)

type Authorizer struct {
	logger *mlg.ZapLog
	config *cnf.Config
}

// Initialize инициализирует синглтон авторизовывальщика с секретным ключом.
func Initialize(c *cnf.Config, l *mlg.ZapLog) (Authorizer, error) {
	au := Authorizer{
		config: c,
		logger: l,
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
		routePattern := chi.RouteContext(r.Context())
		if routePattern.RouteMethod == "POST" && routePattern.URLParams.Values[0] == "orders" {
			// Проверяем формат запроса.
			contentType := r.Header.Get("Content-Type")
			textPlain := strings.Contains(contentType, "text/plain")
			if !textPlain {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		userID, err := au.GetUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		au.logger.ZL.Debug("", zap.Int("userID", userID))
		ctx := context.WithValue(r.Context(), KeyUserIDCtx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (au *Authorizer) SetNewCookie(w http.ResponseWriter, userID int) (err error) {
	au.logger.ZL.Debug("setNewCookie got userID", zap.Int("userID", userID))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Когда создан токен.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(au.config.TokenExp)),
		},
		// Собственное утверждение.
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(au.config.SecretKey))
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
func (au *Authorizer) GetUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, fmt.Errorf("no cookies by name token %w", err)
	}
	// Создаем экземпляр структуры с утверждениями.
	claims := &Claims{}
	// Парсим из строки токена tokenString в структуру claims.
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(au.config.SecretKey), nil
	})
	if err != nil {
		au.logger.ZL.Info("Failed in case to get ownerId from token ", zap.Error(err))
	}
	return claims.UserID, nil
}
