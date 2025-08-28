package order_service_test

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/enums/order_status"
	"loopit/internal/mock"
	"loopit/internal/models"
	"loopit/internal/services/order_service"

	"loopit/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepo(ctrl)

	svc := order_service.NewOrderService(
		mockOrderRepo, nil, nil, logger.NewFakeLogger(),
	)

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError string
	}{
		{
			name: "order not found",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(nil, errors.New("not found"))
			},
			expectedError: "not found",
		},
		{
			name: "invalid status transition",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, Status: order_status.ReturnRequested,
				}, nil)
			},
			expectedError: "order must be in return_requested status",
		},
		{
			name: "repo update error",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, Status: order_status.ReturnRequested,
				}, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(1, order_status.Returned.String()).Return(errors.New("db error"))
			},
			expectedError: "db error",
		},
		{
			name: "success",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, Status: order_status.ReturnRequested,
				}, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(1, order_status.Returned.String()).Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := svc.UpdateOrderStatus(1, order_status.Returned)
			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
				t.Fatalf("expected %q, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetLenderOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepo(ctrl)

	svc := order_service.NewOrderService(mockOrderRepo, nil, nil, logger.NewFakeLogger())

	tests := []struct {
		name          string
		userCtx       *models.UserContext
		setupMocks    func()
		expectErr     bool
		expectedCount int
	}{
		{
			name:       "unauthorized non-lender",
			userCtx:    &models.UserContext{ID: 2, Role: enums.RoleUser},
			setupMocks: func() {},
			expectErr:  true,
		},
		{
			name:    "repo error",
			userCtx: &models.UserContext{ID: 2, Role: enums.RoleLender},
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetLenderOrders(2).Return(nil, errors.New("db fail"))
			},
			expectErr: true,
		},
		{
			name:    "success",
			userCtx: &models.UserContext{ID: 2, Role: enums.RoleLender},
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetLenderOrders(2).Return([]*models.Order{
					{ID: 1}, {ID: 2},
				}, nil)
			},
			expectErr:     false,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			orders, err := svc.GetLenderOrders(tt.userCtx)
			if tt.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectErr && len(orders) != tt.expectedCount {
				t.Fatalf("expected %d orders, got %d", tt.expectedCount, len(orders))
			}
		})
	}
}
