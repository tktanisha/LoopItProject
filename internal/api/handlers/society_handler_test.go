package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"loopit/internal/api/handlers"
	"loopit/internal/constants"
	"loopit/internal/models"
	"loopit/pkg/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mocks the society service interface
type FakeSocietyService struct {
	GetAllFn    func() ([]models.Society, error)
	CreateSocFn func(name, location, pincode string) error
}

func (f *FakeSocietyService) GetAllSocieties() ([]models.Society, error) {
	return f.GetAllFn()
}
func (f *FakeSocietyService) CreateSociety(name, location, pincode string) error {
	return f.CreateSocFn(name, location, pincode)
}

func TestGetAllSocieties(t *testing.T) {
	log := logger.NewFakeLogger()

	tests := []struct {
		name       string
		serviceFn  func() ([]models.Society, error)
		wantStatus int
		wantBody   string
	}{
		{
			name: "service error",
			serviceFn: func() ([]models.Society, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"status":false`,
		},
		{
			name: "success",
			serviceFn: func() ([]models.Society, error) {
				return []models.Society{{ID: 1, Name: "TestSociety"}}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeSocietyService{
				GetAllFn:    tt.serviceFn,
				CreateSocFn: nil,
			}
			h := handlers.NewSocietyHandler(svc, log)

			req := httptest.NewRequest(http.MethodGet, "/societies", nil)
			w := httptest.NewRecorder()

			h.GetAllSocieties(w, req)
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

func TestCreateSociety(t *testing.T) {
	log := logger.NewFakeLogger()

	tests := []struct {
		name       string
		userCtx    interface{}
		body       string
		serviceFn  func(name, location, pincode string) error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "no user context",
			userCtx:    nil,
			body:       `{"name":"A","location":"City","pincode":"12345"}`,
			serviceFn:  nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:       "invalid user context type",
			userCtx:    "not-a-user",
			body:       `{"name":"A","location":"City","pincode":"12345"}`,
			serviceFn:  nil,
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal error",
		},
		{
			name:       "invalid payload",
			userCtx:    &models.UserContext{ID: 1},
			body:       `{invalid-json}`,
			serviceFn:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name:    "service error",
			userCtx: &models.UserContext{ID: 1},
			body:    `{"name":"A","location":"City","pincode":"12345"}`,
			serviceFn: func(name, location, pincode string) error {
				return errors.New("create error")
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "failed to create society",
		},
		{
			name:    "success",
			userCtx: &models.UserContext{ID: 1},
			body:    `{"name":"A","location":"City","pincode":"12345"}`,
			serviceFn: func(name, location, pincode string) error {
				return nil
			},
			wantStatus: http.StatusCreated,
			wantBody:   `"status":true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeSocietyService{
				GetAllFn:    nil,
				CreateSocFn: tt.serviceFn,
			}
			h := handlers.NewSocietyHandler(svc, log)

			req := httptest.NewRequest(http.MethodPost, "/societies", bytes.NewReader([]byte(tt.body)))
			if tt.userCtx != nil {
				ctx := context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			h.CreateSociety(w, req)
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
