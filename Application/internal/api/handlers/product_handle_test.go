package handlers_test

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
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

// // FakeProductService mocks ProductServiceInterface
// type FakeProductService struct {
// 	GetAllFn func() ([]*models.ProductResponse, error)
// 	GetOneFn func(id int) (*models.ProductResponse, error)
// 	CreateFn func(p *models.Product, user *models.UserContext) error
// }

// func (f *FakeProductService) GetAllProducts() ([]*models.ProductResponse, error) {
// 	if f.GetAllFn != nil {
// 		return f.GetAllFn()
// 	}
// 	return nil, nil
// }
// func (f *FakeProductService) GetProductByID(id int) (*models.ProductResponse, error) {
// 	if f.GetOneFn != nil {
// 		return f.GetOneFn(id)
// 	}
// 	return nil, nil
// }
// func (f *FakeProductService) CreateProduct(p *models.Product, user *models.UserContext) error {
// 	if f.CreateFn != nil {
// 		return f.CreateFn(p, user)
// 	}
// 	return nil
// }

// // ----------------- Test GetAllProducts -----------------
// func TestGetAllProducts(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		serviceFn  func() ([]*models.ProductResponse, error)
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{"service error",
// 			func() ([]*models.ProductResponse, error) { return nil, errors.New("db error") },
// 			http.StatusInternalServerError, "failed to fetch products"},
// 		{"success",
// 			func() ([]*models.ProductResponse, error) {
// 				return []*models.ProductResponse{
// 					&models.ProductResponse{Product: models.Product{ID: 1, Name: "Item1"}},
// 				}, nil
// 			},
// 			http.StatusOK, `"status":true`},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeProductService{GetAllFn: tt.serviceFn}
// 			h := handlers.NewProductHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodGet, "/products", nil)
// 			w := httptest.NewRecorder()
// 			h.GetAllProducts(w, req)

// 			if w.Result().StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }

// // ----------------- Test GetProductByID -----------------
// func TestGetProductByID(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		idParam    string
// 		serviceFn  func(id int) (*models.ProductResponse, error)
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{"invalid id", "bad", nil, http.StatusBadRequest, "invalid product id"},
// 		{"service error", "5",
// 			func(id int) (*models.ProductResponse, error) { return nil, errors.New("not found") },
// 			http.StatusNotFound, "product not found"},
// 		{"success", "2",
// 			func(id int) (*models.ProductResponse, error) {
// 				return &models.ProductResponse{Product: models.Product{ID: 2, Name: "Pen"}}, nil
// 			},
// 			http.StatusOK, `"status":true`},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeProductService{GetOneFn: tt.serviceFn}
// 			h := handlers.NewProductHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.idParam, nil)
// 			req.SetPathValue("id", tt.idParam)
// 			w := httptest.NewRecorder()
// 			h.GetProductByID(w, req)

// 			if w.Result().StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }

// // ----------------- Test CreateProduct -----------------
// func TestCreateProduct(t *testing.T) {
// 	log := logger.NewFakeLogger()
// 	productData := models.Product{ID: 10, Name: "Book"}
// 	payload, _ := json.Marshal(productData)

// 	tests := []struct {
// 		name       string
// 		userCtx    any
// 		body       string
// 		serviceFn  func(p *models.Product, u *models.UserContext) error
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{"no user ctx", nil, string(payload), nil, http.StatusUnauthorized, "unauthorized"},
// 		{"invalid JSON", &models.UserContext{ID: 1}, "{bad}", nil, http.StatusBadRequest, "invalid request payload"},
// 		{"service error", &models.UserContext{ID: 1}, string(payload),
// 			func(p *models.Product, u *models.UserContext) error { return errors.New("create fail") },
// 			http.StatusForbidden, "failed to create product"},
// 		{"success", &models.UserContext{ID: 2}, string(payload),
// 			func(p *models.Product, u *models.UserContext) error { return nil },
// 			http.StatusOK, `"status":true`},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeProductService{CreateFn: tt.serviceFn}
// 			h := handlers.NewProductHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodPost, "/products/create", bytes.NewReader([]byte(tt.body)))
// 			if tt.userCtx != nil {
// 				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
// 			}
// 			w := httptest.NewRecorder()
// 			h.CreateProduct(w, req)

// 			if w.Result().StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }
