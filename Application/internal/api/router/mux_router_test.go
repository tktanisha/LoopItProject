package router_test

import (
	"loopit/internal/api/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMuxRouter_HandleFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new ServeMux
	mux := http.NewServeMux()
	muxRouter := router.NewMuxRouter(mux)

	// Define a handler function
	called := false
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	// Register the handler
	muxRouter.HandleFunc("/test", testHandler)

	// Simulate a request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if !called {
		t.Errorf("expected handler to be called, but it was not")
	}
}
