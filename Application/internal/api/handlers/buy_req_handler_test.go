package handlers_test

// import (
// 	"bytes"
// 	"context"
// 	"errors"
// 	"loopit/internal/api/handlers"
// 	"loopit/internal/constants"
// 	"loopit/internal/enums/buyer_request_status"
// 	"loopit/internal/models"
// 	"loopit/pkg/logger"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// // FakeBuyerRequestService mocks the service interface
// type FakeBuyerRequestService struct {
// 	CreateFn func(productID int, user *models.UserContext) error
// 	GetAllFn func(productID int, status buyer_request_status.Status) ([]models.BuyingRequest, error)
// 	UpdateFn func(reqID int, status buyer_request_status.Status, user *models.UserContext) error
// }

// func (f *FakeBuyerRequestService) CreateBuyerRequest(productID int, user *models.UserContext) error {
// 	if f.CreateFn != nil {
// 		return f.CreateFn(productID, user)
// 	}
// 	return nil
// }

// func (f *FakeBuyerRequestService) GetAllBuyerRequestsByStatus(productID int, status buyer_request_status.Status) ([]models.BuyingRequest, error) {
// 	if f.GetAllFn != nil {
// 		return f.GetAllFn(productID, status)
// 	}
// 	return nil, nil
// }

// func (f *FakeBuyerRequestService) UpdateBuyerRequestStatus(reqID int, status buyer_request_status.Status, user *models.UserContext) error {
// 	if f.UpdateFn != nil {
// 		return f.UpdateFn(reqID, status, user)
// 	}
// 	return nil
// }

// // --------------------- Test CreateBuyerRequest ---------------------
// func TestCreateBuyerRequest(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		userCtx    any
// 		body       string
// 		serviceFn  func(productID int, user *models.UserContext) error
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{
// 			name:       "no user context",
// 			userCtx:    nil,
// 			body:       `{"product_id": 1}`,
// 			wantStatus: http.StatusUnauthorized,
// 			wantBody:   "unauthorized",
// 		},
// 		{
// 			name:       "invalid JSON",
// 			userCtx:    &models.UserContext{ID: 1},
// 			body:       `{bad}`,
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "invalid request payload",
// 		},
// 		{
// 			name:    "service error",
// 			userCtx: &models.UserContext{ID: 1},
// 			body:    `{"product_id":123}`,
// 			serviceFn: func(productID int, user *models.UserContext) error {
// 				return errors.New("db error")
// 			},
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "failed to create buyer request",
// 		},
// 		{
// 			name:    "success",
// 			userCtx: &models.UserContext{ID: 2},
// 			body:    `{"product_id":123}`,
// 			serviceFn: func(productID int, user *models.UserContext) error {
// 				return nil
// 			},
// 			wantStatus: http.StatusOK,
// 			wantBody:   `"status":true`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeBuyerRequestService{CreateFn: tt.serviceFn}
// 			h := handlers.NewBuyerRequestHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodPost, "/buyer-requests", bytes.NewReader([]byte(tt.body)))
// 			if tt.userCtx != nil {
// 				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
// 			}
// 			w := httptest.NewRecorder()

// 			h.CreateBuyerRequest(w, req)

// 			resp := w.Result()
// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d, got %d", tt.wantStatus, resp.StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body to contain %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }

// // --------------------- Test GetAllBuyerRequests ---------------------
// func TestGetAllBuyerRequests(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		query      string
// 		serviceFn  func(productID int, status buyer_request_status.Status) ([]models.BuyingRequest, error)
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{
// 			name:       "missing product_id",
// 			query:      "",
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "missing product_id",
// 		},
// 		{
// 			name:       "invalid product_id",
// 			query:      "?product_id=abc",
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "invalid product_id",
// 		},
// 		{
// 			name:       "invalid status",
// 			query:      "?product_id=1&status=wrongStatus",
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "invalid status",
// 		},
// 		{
// 			name:  "service error",
// 			query: "?product_id=1&status=Pending",
// 			serviceFn: func(productID int, status buyer_request_status.Status) ([]models.BuyingRequest, error) {
// 				return nil, errors.New("db error")
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 			wantBody:   "failed to fetch buyer requests",
// 		},
// 		{
// 			name:  "success",
// 			query: "?product_id=1&status=Pending",
// 			serviceFn: func(productID int, status buyer_request_status.Status) ([]models.BuyingRequest, error) {
// 				return []models.BuyingRequest{{ID: 1, ProductID: 1}}, nil
// 			},
// 			wantStatus: http.StatusOK,
// 			wantBody:   `"status":true`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeBuyerRequestService{GetAllFn: tt.serviceFn}
// 			h := handlers.NewBuyerRequestHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodGet, "/buyer-requests"+tt.query, nil)
// 			w := httptest.NewRecorder()
// 			h.GetAllBuyerRequests(w, req)

// 			resp := w.Result()
// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, resp.StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body to contain %q, got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }

// // --------------------- Test UpdateBuyerRequestStatus ---------------------
// func TestUpdateBuyerRequestStatus(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		userCtx    any
// 		reqID      string
// 		body       string
// 		serviceFn  func(reqID int, status buyer_request_status.Status, u *models.UserContext) error
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{
// 			name:       "no user ctx",
// 			userCtx:    nil,
// 			reqID:      "1",
// 			body:       `{"status":"Approved"}`,
// 			wantStatus: http.StatusUnauthorized,
// 			wantBody:   "unauthorized",
// 		},
// 		{
// 			name:       "invalid reqID",
// 			userCtx:    &models.UserContext{ID: 1},
// 			reqID:      "abc",
// 			body:       `{"status":"Pending"}`,
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "invalid buyer request ID",
// 		},
// 		{
// 			name:       "invalid JSON",
// 			userCtx:    &models.UserContext{ID: 1},
// 			reqID:      "1",
// 			body:       `{bad}`,
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "invalid request payload",
// 		},
// 		{
// 			name:       "invalid status",
// 			userCtx:    &models.UserContext{ID: 1},
// 			reqID:      "1",
// 			body:       `{"status":"wrong"}`,
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "invalid status value",
// 		},
// 		{
// 			name:    "service error",
// 			userCtx: &models.UserContext{ID: 1},
// 			reqID:   "1",
// 			body:    `{"status":"Approved"}`,
// 			serviceFn: func(reqID int, status buyer_request_status.Status, u *models.UserContext) error {
// 				return errors.New("update err")
// 			},
// 			wantStatus: http.StatusBadRequest,
// 			wantBody:   "failed to update status",
// 		},
// 		{
// 			name:    "success",
// 			userCtx: &models.UserContext{ID: 1},
// 			reqID:   "1",
// 			body:    `{"status":"Approved"}`,
// 			serviceFn: func(reqID int, status buyer_request_status.Status, u *models.UserContext) error {
// 				return nil
// 			},
// 			wantStatus: http.StatusOK,
// 			wantBody:   `"status":true`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			svc := &FakeBuyerRequestService{UpdateFn: tt.serviceFn}
// 			h := handlers.NewBuyerRequestHandler(svc, log)

// 			req := httptest.NewRequest(http.MethodPatch, "/buyer-requests/"+tt.reqID+"/update", bytes.NewReader([]byte(tt.body)))
// 			if tt.userCtx != nil {
// 				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
// 			}
// 			// set path param manually since httptest doesn't parse it
// 			req.SetPathValue("requestId", tt.reqID)

// 			w := httptest.NewRecorder()
// 			h.UpdateBuyerRequestStatus(w, req)

// 			resp := w.Result()
// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d, got %d", tt.wantStatus, resp.StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body to contain %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }
