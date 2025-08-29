package handlers_test

import (
	"context"
	"errors"

	"loopit/internal/api/handlers"
	"loopit/internal/constants"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/pkg/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeService struct {
	BecomeLenderFn func(*models.UserContext) error
}

func (f *FakeService) BecomeLender(u *models.UserContext) error {
	return f.BecomeLenderFn(u)
}

func TestBecomeLenderHandler(t *testing.T) {
	var log logger.LoggerInterface = logger.NewFakeLogger()

	tests := []struct {
		name       string
		userCtx    interface{}
		serviceFn  func(*models.UserContext) error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "no user context",
			userCtx:    nil,
			serviceFn:  nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:       "invalid user context type",
			userCtx:    "not-a-user",
			serviceFn:  nil,
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal error",
		},
		{
			name:    "service error",
			userCtx: &models.UserContext{ID: 1, Role: enums.RoleUser},
			serviceFn: func(u *models.UserContext) error {
				return errors.New("failed")
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":false`,
		},
		{
			name:    "success",
			userCtx: &models.UserContext{ID: 2, Role: enums.RoleUser},
			serviceFn: func(u *models.UserContext) error {
				return nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeService{BecomeLenderFn: tt.serviceFn}
			h := handlers.NewUserHandler(svc, log)

			req := httptest.NewRequest(http.MethodPatch, "/users/become-lender", strings.NewReader(""))
			if tt.userCtx != nil {
				ctx := context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			h.BecomeLender(w, req)
			resp := w.Result()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			body := w.Body.String()

			if !strings.Contains(body, tt.wantBody) {
				t.Errorf("expected body to contain %q, got %s", tt.wantBody, body)
			}
		})
	}
}
