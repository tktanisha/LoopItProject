package middleware

import (
	"context"
	"loopit/internal/constants"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	"net/http"
	"strings"
)

// AuthMiddleware provides JWT authentication for protected routes
func AuthMiddleware(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		role, err := enums.ParseRole(claims.Role)
		if err != nil {
			http.Error(w, "invalid role", http.StatusUnauthorized)
			return
		}

		userCtx := &models.UserContext{
			ID:   claims.UserID,
			Role: role,
		}

		ctx := context.WithValue(r.Context(), constants.UserCtxKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
