package api_test

import (
	"testing"

	"loopit/internal/api"
	"loopit/internal/mock"

	"github.com/golang/mock/gomock"
)

func TestSetupRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock router
	r := mock.NewMockRouter(ctrl)

	// Mock handlers
	h1 := mock.NewMockHandler(ctrl)
	h2 := mock.NewMockHandler(ctrl)

	// Expectations: Each handler should call RegisterRoutes once with the same router
	h1.EXPECT().RegisterRoutes(r).Times(1)
	h2.EXPECT().RegisterRoutes(r).Times(1)

	// Call the function under test
	api.SetupRoutes(r, h1, h2)
}
