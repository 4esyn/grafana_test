package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"strings"
)

func AuthMiddleware(tokenAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}

			// Удаляем префикс "Bearer " если он есть
			token = strings.TrimPrefix(token, "Bearer ")

			// Проверяем токен
			_, err := jwtauth.VerifyRequest(tokenAuth, r, jwtauth.TokenFromHeader)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
