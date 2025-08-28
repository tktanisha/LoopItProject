package handlers_test

import (
	"bytes"
	"errors"
	"loopit/internal/api/handlers"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/pkg/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// FakeAuthService mocks AuthServiceInterface
type FakeAuthService struct {
	LoginFn    func(email, password string) (string, *models.User, error)
	RegisterFn func(user *models.User) error
}

func (f *FakeAuthService) Login(email, password string) (string, *models.User, error) {
	if f.LoginFn != nil {
		return f.LoginFn(email, password)
	}
	return "", nil, nil
}
func (f *FakeAuthService) Register(user *models.User) error {
	if f.RegisterFn != nil {
		return f.RegisterFn(user)
	}
	return nil
}

// -------------------- Test Login --------------------
func TestLogin(t *testing.T) {
	log := logger.NewFakeLogger()

	tests := []struct {
		name       string
		body       string
		serviceFn  func(email, password string) (string, *models.User, error)
		wantStatus int
		wantBody   string
	}{
		{
			name:       "invalid JSON",
			body:       `{bad-json}`,
			serviceFn:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name: "auth error",
			body: `{"email":"test@example.com","password":"wrong"}`,
			serviceFn: func(email, password string) (string, *models.User, error) {
				return "", nil, errors.New("invalid creds")
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "invalid credentials",
		},
		{
			name: "success",
			body: `{"email":"test@example.com","password":"pass"}`,
			serviceFn: func(email, password string) (string, *models.User, error) {
				return "fake-token", &models.User{ID: 1, FullName: "John Doe", Role: enums.RoleUser}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `"token":"fake-token"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeAuthService{LoginFn: tt.serviceFn}
			h := handlers.NewAuthHandler(svc, log)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(tt.body)))
			w := httptest.NewRecorder()
			h.Login(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected %d got %d", tt.wantStatus, resp.StatusCode)
			}
			if !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("expected body to contain %q got %s", tt.wantBody, w.Body.String())
			}
		})
	}
}

// -------------------- Test Register --------------------
func TestRegister(t *testing.T) {
	log := logger.NewFakeLogger()
	registerPayload := `{"fullname":"John","email":"test@example.com","password":"pass","phone_number":"12345","address":"Earth"}`

	tests := []struct {
		name       string
		body       string
		registerFn func(*models.User) error
		loginFn    func(string, string) (string, *models.User, error)
		wantStatus int
		wantBody   string
	}{
		{
			name:       "invalid JSON",
			body:       `{bad-json}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name: "registration error",
			body: registerPayload,
			registerFn: func(u *models.User) error {
				return errors.New("reg error")
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "registration failed",
		},
		{
			name:       "auto login failed",
			body:       registerPayload,
			registerFn: func(u *models.User) error { return nil },
			loginFn: func(e, p string) (string, *models.User, error) {
				return "", nil, errors.New("login fail")
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "invalid credentials",
		},
		{
			name:       "success",
			body:       registerPayload,
			registerFn: func(u *models.User) error { return nil },
			loginFn: func(e, p string) (string, *models.User, error) {
				return "token-123", &models.User{ID: 1, FullName: "John", Role: enums.RoleUser}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `"token":"token-123"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeAuthService{
				RegisterFn: tt.registerFn,
				LoginFn:    tt.loginFn,
			}
			h := handlers.NewAuthHandler(svc, log)

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(tt.body)))
			w := httptest.NewRecorder()
			h.Register(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected %d got %d", tt.wantStatus, resp.StatusCode)
			}
			if !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("expected body to contain %q got %s", tt.wantBody, w.Body.String())
			}
		})
	}
}

// -------------------- Test Logout --------------------
func TestLogout(t *testing.T) {
	log := logger.NewFakeLogger()
	h := handlers.NewAuthHandler(&FakeAuthService{}, log)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()
	h.Logout(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"status":true`) {
		t.Errorf("expected status true got %s", body)
	}
	if !strings.Contains(body, "logged out successfully") {
		t.Errorf("expected logout message got %s", body)
	}
}
