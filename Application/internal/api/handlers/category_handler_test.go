package handlers_test

// import (
// 	"bytes"
// 	"context"
// 	"errors"
// 	"loopit/internal/api/handlers"
// 	"loopit/internal/constants"
// 	"loopit/internal/models"
// 	"loopit/pkg/logger"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// // FakeCategoryService mocks the category service interface
// type FakeCategoryService struct {
// 	CreateFn func(name string, price, security float64) error
// 	GetAllFn func() ([]models.Category, error)
// }

// func (f *FakeCategoryService) CreateCategory(name string, price, security float64) error {
// 	if f.CreateFn == nil {
// 		return nil
// 	}
// 	return f.CreateFn(name, price, security)
// }

// func (f *FakeCategoryService) GetAllCategories() ([]models.Category, error) {
// 	if f.GetAllFn == nil {
// 		return nil, nil
// 	}
// 	return f.GetAllFn()
// }

// // --------------------- Test GetAllCategories ---------------------
// func TestGetAllCategories(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		serviceFn  func() ([]models.Category, error)
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{
// 			name: "service error",
// 			serviceFn: func() ([]models.Category, error) {
// 				return nil, errors.New("db error")
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 			wantBody:   `"status":false`,
// 		},
// 		{
// 			name: "success",
// 			serviceFn: func() ([]models.Category, error) {
// 				return []models.Category{{ID: 1, Name: "Electronics"}}, nil
// 			},
// 			wantStatus: http.StatusOK,
// 			wantBody:   `"status":true`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeCategoryService{GetAllFn: tt.serviceFn}
// 			h := handlers.NewCategoryHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
// 			w := httptest.NewRecorder()

// 			h.GetAllCategories(w, req)
// 			resp := w.Result()

// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d, got %d", tt.wantStatus, resp.StatusCode)
// 			}
// 			body := w.Body.String()
// 			if !strings.Contains(body, tt.wantBody) {
// 				t.Errorf("expected body to contain %q, got %s", tt.wantBody, body)
// 			}
// 		})
// 	}
// }

// // --------------------- Test CreateCategory ---------------------
// func TestCreateCategory(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		userCtx    any
// 		body       string
// 		serviceFn  func(name string, price, security float64) error
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{
// 			name:       "no user context",
// 			userCtx:    nil,
// 			body:       `{"name":"Laptop","price":1000,"security":500}`,
// 			serviceFn:  nil,
// 			wantStatus: http.StatusUnauthorized,
// 			wantBody:   "unauthorized",
// 		},
// 		{
// 			name:       "invalid user context type",
// 			userCtx:    "not-a-user",
// 			body:       `{"name":"Laptop","price":1000,"security":500}`,
// 			serviceFn:  nil,
// 			wantStatus: http.StatusInternalServerError,
// 			wantBody:   "internal error",
// 		},
// 		{
// 			name:       "invalid JSON payload",
// 			userCtx:    &models.UserContext{ID: 1},
// 			body:       `{bad-json}`,
// 			serviceFn:  nil,
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "invalid request payload",
// 		},
// 		{
// 			name:    "service error",
// 			userCtx: &models.UserContext{ID: 1},
// 			body:    `{"name":"Laptop","price":1000,"security":500}`,
// 			serviceFn: func(name string, price, security float64) error {
// 				return errors.New("create error")
// 			},
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "failed to create category",
// 		},
// 		{
// 			name:    "success",
// 			userCtx: &models.UserContext{ID: 1},
// 			body:    `{"name":"Laptop","price":1000,"security":500}`,
// 			serviceFn: func(name string, price, security float64) error {
// 				return nil
// 			},
// 			wantStatus: http.StatusCreated,
// 			wantBody:   `"status":true`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeCategoryService{CreateFn: tt.serviceFn}
// 			h := handlers.NewCategoryHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte(tt.body)))
// 			if tt.userCtx != nil {
// 				ctx := context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx)
// 				req = req.WithContext(ctx)
// 			}
// 			w := httptest.NewRecorder()

// 			h.CreateCategory(w, req)
// 			resp := w.Result()

// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d, got %d", tt.wantStatus, resp.StatusCode)
// 			}
// 			body := w.Body.String()
// 			if !strings.Contains(body, tt.wantBody) {
// 				t.Errorf("expected body to contain %q, got %s", tt.wantBody, body)
// 			}
// 		})
// 	}
// }
