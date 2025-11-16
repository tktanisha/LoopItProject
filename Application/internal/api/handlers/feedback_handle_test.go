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

// FakeFeedbackService mocks feedback service interface
type FakeFeedbackService struct {
	GiveFn     func(orderID int, feedbackText string, rating int, user *models.UserContext) error
	GivenFn    func(user *models.UserContext) ([]models.Feedback, error)
	ReceivedFn func(user *models.UserContext) ([]models.Feedback, error)
}

func (f *FakeFeedbackService) GiveFeedback(orderID int, feedbackText string, rating int, user *models.UserContext) error {
	if f.GiveFn != nil {
		return f.GiveFn(orderID, feedbackText, rating, user)
	}
	return nil
}
func (f *FakeFeedbackService) GetAllGivenFeedbacks(user *models.UserContext) ([]models.Feedback, error) {
	if f.GivenFn != nil {
		return f.GivenFn(user)
	}
	return nil, nil
}
func (f *FakeFeedbackService) GetAllReceivedFeedbacks(user *models.UserContext) ([]models.Feedback, error) {
	if f.ReceivedFn != nil {
		return f.ReceivedFn(user)
	}
	return nil, nil
}

// ----------------- Test GiveFeedback -----------------
func TestGiveFeedback(t *testing.T) {
	log := logger.NewFakeLogger()

	tests := []struct {
		name       string
		userCtx    any
		body       string
		serviceFn  func(orderID int, feedbackText string, rating int, u *models.UserContext) error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "no user context",
			userCtx:    nil,
			body:       `{"order_id":1,"feedback_text":"good","rating":5}`,
			serviceFn:  nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:       "invalid user context type",
			userCtx:    "not-a-user",
			body:       `{"order_id":1,"feedback_text":"good","rating":5}`,
			serviceFn:  nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name:       "invalid JSON payload",
			userCtx:    &models.UserContext{ID: 1},
			body:       `{bad-json}`,
			serviceFn:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name:    "service error",
			userCtx: &models.UserContext{ID: 1},
			body:    `{"order_id":1,"feedback_text":"good","rating":5}`,
			serviceFn: func(orderID int, feedbackText string, rating int, u *models.UserContext) error {
				return errors.New("save error")
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "failed to give feedback",
		},
		{
			name:    "success",
			userCtx: &models.UserContext{ID: 1},
			body:    `{"order_id":1,"feedback_text":"good","rating":5}`,
			serviceFn: func(orderID int, feedbackText string, rating int, u *models.UserContext) error {
				return nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeFeedbackService{GiveFn: tt.serviceFn}
			h := handlers.NewFeedbackHandler(svc, log)

			req := httptest.NewRequest(http.MethodPost, "/feedbacks", bytes.NewReader([]byte(tt.body)))
			if tt.userCtx != nil {
				ctx := context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			h.GiveFeedback(w, req)

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

// ----------------- Test GetAllGivenFeedbacks -----------------
func TestGetAllGivenFeedbacks(t *testing.T) {
	log := logger.NewFakeLogger()

	tests := []struct {
		name       string
		userCtx    any
		serviceFn  func(u *models.UserContext) ([]models.Feedback, error)
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
			name:    "service error",
			userCtx: &models.UserContext{ID: 1},
			serviceFn: func(u *models.UserContext) ([]models.Feedback, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed to fetch given feedbacks",
		},
		{
			name:    "success",
			userCtx: &models.UserContext{ID: 1},
			serviceFn: func(u *models.UserContext) ([]models.Feedback, error) {
				return []models.Feedback{{ID: 1, Text: "Great service!"}}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeFeedbackService{GivenFn: tt.serviceFn}
			h := handlers.NewFeedbackHandler(svc, log)

			req := httptest.NewRequest(http.MethodGet, "/feedbacks/given", nil)
			if tt.userCtx != nil {
				ctx := context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			h.GetAllGivenFeedbacks(w, req)

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

// ----------------- Test GetAllReceivedFeedbacks -----------------
func TestGetAllReceivedFeedbacks(t *testing.T) {
	log := logger.NewFakeLogger()

	tests := []struct {
		name       string
		userCtx    any
		serviceFn  func(u *models.UserContext) ([]models.Feedback, error)
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
			name:    "service error",
			userCtx: &models.UserContext{ID: 1},
			serviceFn: func(u *models.UserContext) ([]models.Feedback, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed to fetch received feedbacks",
		},
		{
			name:    "success",
			userCtx: &models.UserContext{ID: 1},
			serviceFn: func(u *models.UserContext) ([]models.Feedback, error) {
				return []models.Feedback{{ID: 2, Text: "Excellent!"}}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeFeedbackService{ReceivedFn: tt.serviceFn}
			h := handlers.NewFeedbackHandler(svc, log)

			req := httptest.NewRequest(http.MethodGet, "/feedbacks/received", nil)
			if tt.userCtx != nil {
				ctx := context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			h.GetAllReceivedFeedbacks(w, req)

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
