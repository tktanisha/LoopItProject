package initializer_test

import (
	"testing"

	"loopit/internal/initializer"
	"loopit/internal/mock"

	"github.com/golang/mock/gomock"
)

// TestInitDBRepos ensures that InitDBRepos initializes all repositories correctly.
func TestInitDBRepos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockLogger := mock.NewMockLoggerInterface(ctrl)
	mockDB := mock.NewMockDatabaseInterface(ctrl)

	// Call the function under test
	err := initializer.InitDBRepos(mockLogger, mockDB)
	if err != nil {
		t.Fatalf("InitDBRepos() returned unexpected error: %v", err)
	}

	// Validate that repositories are initialized
	if initializer.LenderRepo == nil {
		t.Error("LenderRepo was not initialized")
	}
	if initializer.UserRepo == nil {
		t.Error("UserRepo was not initialized")
	}
	if initializer.CategoryRepo == nil {
		t.Error("CategoryRepo was not initialized")
	}
	if initializer.ProductRepo == nil {
		t.Error("ProductRepo was not initialized")
	}
	if initializer.BuyerRequestRepo == nil {
		t.Error("BuyerRequestRepo was not initialized")
	}
	if initializer.OrderRepo == nil {
		t.Error("OrderRepo was not initialized")
	}
	if initializer.ReturnRequestRepo == nil {
		t.Error("ReturnRequestRepo was not initialized")
	}
	if initializer.FeedBackRepo == nil {
		t.Error("FeedBackRepo was not initialized")
	}
	if initializer.SocietyRepo == nil {
		t.Error("SocietyRepo was not initialized")
	}
}
