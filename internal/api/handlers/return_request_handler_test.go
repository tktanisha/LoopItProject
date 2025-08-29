package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"loopit/internal/api/handlers"
	"loopit/internal/constants"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/pkg/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// FakeReturnRequestService mocks ReturnRequestServiceInterface
type FakeReturnRequestService struct {
	GetPendingFn func(userID int) ([]models.ReturnRequest, error)
	CreateFn     func(userID int, orderID int) error
	UpdateFn     func(userID, reqID int, status return_request_status.Status) error
}

func (f *FakeReturnRequestService) GetPendingReturnRequests(userID int) ([]models.ReturnRequest, error) {
	if f.GetPendingFn != nil {
		return f.GetPendingFn(userID)
	}
	return nil, nil
}
func (f *FakeReturnRequestService) CreateReturnRequest(userID int, orderID int) error {
	if f.CreateFn != nil {
		return f.CreateFn(userID, orderID)
	}
	return nil
}
func (f *FakeReturnRequestService) UpdateReturnRequestStatus(userID, reqID int, status return_request_status.Status) error {
	if f.UpdateFn != nil {
		return f.UpdateFn(userID, reqID, status)
	}
	return nil
}

// -------------------- Test GetPendingReturnRequests -------------------
func TestGetPendingReturnRequests(t *testing.T) {
	log := logger.NewFakeLogger()

	tests := []struct {
		name       string
		userCtx    interface{}
		serviceFn  func(userID int) ([]models.ReturnRequest, error)
		wantStatus int
		wantBody   string
	}{
		{"no user context", nil, nil, http.StatusUnauthorized, "unauthorized"},
		{"service error", &models.UserContext{ID: 1}, func(id int) ([]models.ReturnRequest, error) {
			return nil, errors.New("db error")
		}, http.StatusInternalServerError, "could not fetch return requests"},
		{"success", &models.UserContext{ID: 2}, func(id int) ([]models.ReturnRequest, error) {
			return []models.ReturnRequest{{ID: 1, OrderID: 10}}, nil
		}, http.StatusOK, `"status":true`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeReturnRequestService{GetPendingFn: tt.serviceFn}
			h := handlers.NewReturnRequestHandler(svc, log)

			req := httptest.NewRequest(http.MethodGet, "/return-requests", nil)
			if tt.userCtx != nil {
				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
			}
			w := httptest.NewRecorder()
			h.GetPendingReturnRequests(w, req)

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, w.Result().StatusCode)
			}
			if !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("expected body %q, got %s", tt.wantBody, w.Body.String())
			}
		})
	}
}

// -------------------- Test CreateReturnRequest -------------------
func TestCreateReturnRequest(t *testing.T) {
	log := logger.NewFakeLogger()
	payload := `{"order_id": 15}`

	tests := []struct {
		name       string
		userCtx    interface{}
		body       string
		serviceFn  func(userID int, orderID int) error
		wantStatus int
		wantBody   string
	}{
		{"no user context", nil, payload, nil, http.StatusUnauthorized, "unauthorized"},
		{"invalid JSON", &models.UserContext{ID: 1}, `{bad-json}`, nil, http.StatusBadRequest, "invalid request payload"},
		{"service error", &models.UserContext{ID: 2}, payload,
			func(u, o int) error { return errors.New("fail") },
			http.StatusBadRequest, "could not create return request"},
		{"success", &models.UserContext{ID: 3}, payload,
			func(u, o int) error { return nil },
			http.StatusOK, `"status":true`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeReturnRequestService{CreateFn: tt.serviceFn}
			h := handlers.NewReturnRequestHandler(svc, log)

			req := httptest.NewRequest(http.MethodPost, "/return-requests", bytes.NewReader([]byte(tt.body)))
			if tt.userCtx != nil {
				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
			}
			w := httptest.NewRecorder()
			h.CreateReturnRequest(w, req)

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
			}
			if !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("expected body to contain %q, got %s", tt.wantBody, w.Body.String())
			}
		})
	}
}

// -------------------- Test UpdateReturnRequestStatus -------------------
func TestUpdateReturnRequestStatus(t *testing.T) {
	log := logger.NewFakeLogger()
	validBody := `{"status":"Approved"}`

	tests := []struct {
		name       string
		userCtx    interface{}
		reqID      string
		body       string
		serviceFn  func(userID, reqID int, status return_request_status.Status) error
		wantStatus int
		wantBody   string
	}{
		{"no context", nil, "1", validBody, nil, http.StatusUnauthorized, "unauthorized"},
		{"invalid reqID", &models.UserContext{ID: 1}, "bad", validBody, nil, http.StatusBadRequest, "invalid request id"},
		{"invalid JSON", &models.UserContext{ID: 2}, "1", "{bad}", nil, http.StatusBadRequest, "invalid request payload"},
		{"invalid status", &models.UserContext{ID: 2}, "1", `{"status":"WRONG"}`, nil, http.StatusBadRequest, "invalid status"},
		{"service error", &models.UserContext{ID: 3}, "5", validBody,
			func(u, r int, s return_request_status.Status) error { return errors.New("fail") },
			http.StatusBadRequest, "could not update return request"},
		{"success", &models.UserContext{ID: 4}, "6", validBody,
			func(u, r int, s return_request_status.Status) error { return nil },
			http.StatusOK, `"status":true`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeReturnRequestService{UpdateFn: tt.serviceFn}
			h := handlers.NewReturnRequestHandler(svc, log)

			req := httptest.NewRequest(http.MethodPatch, "/return-requests/"+tt.reqID+"/update", bytes.NewReader([]byte(tt.body)))
			req.SetPathValue("requestId", tt.reqID)
			if tt.userCtx != nil {
				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
			}
			w := httptest.NewRecorder()
			h.UpdateReturnRequestStatus(w, req)

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
			}
			if !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("expected body %q got %s", tt.wantBody, w.Body.String())
			}
		})
	}
}
