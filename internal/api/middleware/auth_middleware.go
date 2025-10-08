package middleware

import (
	"context"
	"fmt"
	"loopit/internal/constants"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	"net/http"
	"strings"
)

// AuthMiddleware provides JWT authentication for protected routes
func AuthMiddleware(log logger.LoggerInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// http.Error(w, "missinzg or invalid authorization header", http.StatusUnauthorized)
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "missing or invalid authorization header", "authorization header must be in format 'Bearer <token>'")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			// http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "invalid or expired token", "token is either invalid or has expired")
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
		fmt.Println("to move in the handler")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
