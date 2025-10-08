package initializer_test

import (
	"testing"

	"loopit/internal/initializer"
	"loopit/pkg/logger"
)

func TestInitServices(t *testing.T) {
	// Backup original services to restore later
	originalAuthService := initializer.AuthService
	originalProductService := initializer.ProductService
	originalCategoryService := initializer.CategoryService
	originalUserService := initializer.UserService
	originalBuyerRequestService := initializer.BuyerRequestService
	originalOrderService := initializer.OrderService
	originalReturnRequestService := initializer.ReturnRequestService
	originalFeedbackService := initializer.FeedBackService
	originalSocietyService := initializer.SocietyService

	defer func() {
		initializer.AuthService = originalAuthService
		initializer.ProductService = originalProductService
		initializer.CategoryService = originalCategoryService
		initializer.UserService = originalUserService
		initializer.BuyerRequestService = originalBuyerRequestService
		initializer.OrderService = originalOrderService
		initializer.ReturnRequestService = originalReturnRequestService
		initializer.FeedBackService = originalFeedbackService
		initializer.SocietyService = originalSocietyService
	}()

	// Call the function
	l := logger.NewFakeLogger()
	initializer.InitServiceFiles(l)

	// Verify that the services are initialized (non-nil check)
	if initializer.AuthService == nil {
		t.Error("expected AuthService to be initialized, got nil")
	}
	if initializer.ProductService == nil {
		t.Error("expected ProductService to be initialized, got nil")
	}
	if initializer.CategoryService == nil {
		t.Error("expected CategoryService to be initialized, got nil")
	}
	if initializer.UserService == nil {
		t.Error("expected UserService to be initialized, got nil")
	}
	if initializer.BuyerRequestService == nil {
		t.Error("expected BuyerRequestService to be initialized, got nil")
	}
	if initializer.OrderService == nil {
		t.Error("expected OrderService to be initialized, got nil")
	}
	if initializer.ReturnRequestService == nil {
		t.Error("expected ReturnRequestService to be initialized, got nil")
	}
	if initializer.FeedBackService == nil {
		t.Error("expected FeedBackService to be initialized, got nil")
	}
	if initializer.SocietyService == nil {
		t.Error("expected SocietyService to be initialized, got nil")
	}
}
