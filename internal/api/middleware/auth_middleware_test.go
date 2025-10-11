package middleware_test

import (
	"fmt"
	"io/ioutil"
	"loopit/internal/api/middleware"
	"loopit/internal/constants"
	"loopit/internal/enums"
	"loopit/pkg/logger"

	"loopit/internal/models"
	"loopit/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		setupToken     func() string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid authorization header format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid jwt token",
			authHeader:     "Bearer junk.token.value",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name: "valid jwt token",
			setupToken: func() string {
				token, _ := utils.GenerateJWT(123, enums.RoleUser.String())
				return token
			},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var calledNext bool
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calledNext = true
				val := r.Context().Value(constants.UserCtxKey)
				if val == nil {
					t.Errorf("expected UserCtx in context, got nil")
				}
				if userCtx, ok := val.(*models.UserContext); ok {
					if userCtx.ID != 123 && tt.expectNext {
						t.Errorf("expected UserID=123, got %d", userCtx.ID)
					}
					if userCtx.Role != enums.RoleUser && tt.expectNext {
						t.Errorf("expected Role=User, got %v", userCtx.Role)
					}
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)

		
			if tt.setupToken != nil {
				token := tt.setupToken()
				tt.authHeader = "Bearer " + token
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := middleware.AuthMiddleware(logger.NewFakeLogger(), nextHandler)

			handler.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)
			fmt.Printf("%s -> %d %s\n", tt.name, res.StatusCode, string(body))

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			if tt.expectNext && !calledNext {
				t.Errorf("expected next handler to be called, but it was not")
			}
			if !tt.expectNext && calledNext {
				t.Errorf("expected next handler NOT to be called, but it was")
			}
		})
	}
}
