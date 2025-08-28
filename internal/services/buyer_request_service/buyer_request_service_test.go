package buyer_request_service_test

import (
	"errors"
	"fmt"
	"loopit/internal/enums"
	br_status "loopit/internal/enums/buyer_request_status"
	"loopit/internal/mock"
	"loopit/internal/models"
	"loopit/internal/services/buyer_request_service"

	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreateBuyerRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuyerRepo := mock.NewMockBuyerRequestRepo(ctrl)
	mockProductRepo := mock.NewMockProductRepo(ctrl)

	// We don’t need categoryRepo and orderRepo for this method
	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
	mockCategoryRepo := mock.NewMockCategoryRepo(ctrl)

	svc := buyer_request_service.NewBuyerRequestService(
		mockBuyerRepo, mockProductRepo, mockOrderRepo, mockCategoryRepo, nil,
	)

	userCtx := &models.UserContext{ID: 2, Role: enums.RoleUser}

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError string
	}{
		{
			name: "product not found",
			setupMocks: func() {
				mockProductRepo.EXPECT().FindByID(1).Return(nil, errors.New("not found"))
			},
			expectedError: "product not found",
		},
		{
			name: "product not available",
			setupMocks: func() {
				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
					Product: models.Product{IsAvailable: false},
				}, nil)
			},
			expectedError: "product not available",
		},
		{
			name: "lender requesting own product",
			setupMocks: func() {
				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
					Product: models.Product{IsAvailable: true, LenderID: 2},
				}, nil)
			},
			expectedError: "lender cannot create a buying request",
		},
		{
			name: "duplicate request exists",
			setupMocks: func() {
				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
					Product: models.Product{IsAvailable: true, LenderID: 10},
				}, nil)
				mockBuyerRepo.EXPECT().GetAllBuyerRequests(gomock.Any()).Return([]models.BuyingRequest{
					{ProductID: 1, RequestedBy: 2, Status: br_status.Pending},
				}, nil)
			},
			expectedError: "a pending or approved request already exists",
		},
		{
			name: "success",
			setupMocks: func() {
				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
					Product: models.Product{IsAvailable: true, LenderID: 10},
				}, nil)
				mockBuyerRepo.EXPECT().GetAllBuyerRequests(gomock.Any()).Return(nil, nil)
				mockBuyerRepo.EXPECT().CreateBuyerRequest(gomock.Any()).Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := svc.CreateBuyerRequest(1, userCtx)
			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}
			if err != nil && tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
				t.Fatalf("expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}

// helper
func contains(got, want string) bool {
	return fmt.Sprintf("%v", got) == want || (len(got) >= len(want) && got[:len(want)] == want)
}
