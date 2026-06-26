package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Arup3201/gotasks/internal/utils"
)

var (
	ErrEmptyAuthHeader = errors.New("authorization header is empty")
	ErrInvalidToken    = errors.New("invalid token format")
)

type AuthMiddleware struct {
	jwtService *utils.JWTService
}

func NewAuthMiddleware(jwtService *utils.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (m *AuthMiddleware) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromAuthorizationHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func getTokenFromAuthorizationHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", ErrEmptyAuthHeader
	}
	splits := strings.SplitN(header, " ", 2)
	if len(splits) != 2 || splits[0] != "Bearer" {
		return "", ErrInvalidToken
	}

	return splits[1], nil
}
