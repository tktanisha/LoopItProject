package middleware

import (
	"log"
	"loopit/pkg/utils"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Fatalf("panic recovered: %v", rec)
				utils.WriteErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred", "")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
