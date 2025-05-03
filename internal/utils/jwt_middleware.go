package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/srgjo27/e-learning/internal/entity"
	"github.com/srgjo27/e-learning/internal/usecase"
)

func JWTMiddleware(authUseCase *usecase.AuthUseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		userID, roleStr, err := authUseCase.ParseToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "userRole", entity.Role(roleStr))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}