package order_service_test

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/enums/order_status"
	"loopit/internal/enums/return_request_status"
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
	svc := order_service.NewOrderService(mockOrderRepo, nil, nil, logger.NewFakeLogger())

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError string
	}{
		{
			name: "repo error",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(nil, errors.New("not found"))
			},
			expectedError: "not found",
		},
		{
			name: "order nil",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(nil, nil)
			},
			expectedError: "order not found",
		},
		{
			name: "invalid transition",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, Status: order_status.Returned,
				}, nil)
			},
			expectedError: "order must be in return_requested status",
		},
		{
			name: "update error",
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

func TestGetOrderHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
	svc := order_service.NewOrderService(mockOrderRepo, nil, nil, logger.NewFakeLogger())

	userCtx := &models.UserContext{ID: 10, Role: enums.RoleUser}
	statuses := []order_status.Status{order_status.InUse, order_status.Returned}

	t.Run("repo error", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetOrderHistory(userCtx.ID, []string{"In Use", "Returned"}).
			Return(nil, errors.New("db fail"))
		_, err := svc.GetOrderHistory(userCtx, statuses)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetOrderHistory(userCtx.ID, []string{"In Use", "Returned"}).
			Return([]*models.Order{{ID: 1}, {ID: 2}}, nil)
		orders, err := svc.GetOrderHistory(userCtx, statuses)
		if err != nil || len(orders) != 2 {
			t.Fatalf("expected 2 orders, got %v (err: %v)", orders, err)
		}
	})
}

func TestMarkOrderAsReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
	mockProductRepo := mock.NewMockProductRepo(ctrl)
	mockReturnRequestRepo := mock.NewMockReturnRequestRepo(ctrl)

	svc := order_service.NewOrderService(mockOrderRepo, mockReturnRequestRepo, mockProductRepo, logger.NewFakeLogger())

	userCtx := &models.UserContext{ID: 5, Role: enums.RoleLender}
	order := &models.Order{ID: 1, ProductID: 99}

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError string
	}{
		{
			name: "order not found",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(nil, errors.New("db error"))
			},
			expectedError: "db error",
		},
		{
			name: "order nil",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(nil, nil)
			},
			expectedError: "order not found",
		},
		{
			name: "product fetch error",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(order, nil)
				mockProductRepo.EXPECT().FindByID(99).Return(nil, errors.New("prod fail"))
			},
			expectedError: "unable to find product",
		},
		{
			name: "product nil",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(order, nil)
				mockProductRepo.EXPECT().FindByID(99).Return(nil, nil)
			},
			expectedError: "product not found",
		},
		{
			name: "unauthorized lender",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(order, nil)
				mockProductRepo.EXPECT().FindByID(99).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 99},
				}, nil)
			},
			expectedError: "unauthorized lender",
		},
		{
			name: "return requests fetch error",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(order, nil)
				mockProductRepo.EXPECT().FindByID(99).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 5},
				}, nil)
				mockReturnRequestRepo.EXPECT().
					GetAllReturnRequests([]string{return_request_status.Approved.String()}).
					Return(nil, errors.New("fail"))
			},
			expectedError: "unable to find return requests",
		},
		{
			name: "no approved request",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(order, nil)
				mockProductRepo.EXPECT().FindByID(99).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 5},
				}, nil)
				mockReturnRequestRepo.EXPECT().
					GetAllReturnRequests([]string{return_request_status.Approved.String()}).
					Return([]models.ReturnRequest{}, nil)
			},
			expectedError: "order has not been approved for return",
		},
		{
			name: "update status error",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(order, nil)
				mockProductRepo.EXPECT().FindByID(99).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 5},
				}, nil)
				mockReturnRequestRepo.EXPECT().
					GetAllReturnRequests([]string{return_request_status.Approved.String()}).
					Return([]models.ReturnRequest{{OrderID: 1}}, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(1, order_status.Returned.String()).
					Return(errors.New("db fail"))
			},
			expectedError: "db fail",
		},
		{
			name: "success",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(order, nil)
				mockProductRepo.EXPECT().FindByID(99).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 5},
				}, nil)
				mockReturnRequestRepo.EXPECT().
					GetAllReturnRequests([]string{return_request_status.Approved.String()}).
					Return([]models.ReturnRequest{{OrderID: 1}}, nil)
				mockOrderRepo.EXPECT().UpdateOrderStatus(1, order_status.Returned.String()).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := svc.MarkOrderAsReturned(1, userCtx)
			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
				t.Fatalf("expected %q, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetAllApprovedAwaitingOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
	mockReturnRequestRepo := mock.NewMockReturnRequestRepo(ctrl)

	svc := order_service.NewOrderService(mockOrderRepo, mockReturnRequestRepo, nil, logger.NewFakeLogger())

	userCtx := &models.UserContext{ID: 5, Role: enums.RoleLender}

	t.Run("unauthorized", func(t *testing.T) {
		_, err := svc.GetAllApprovedAwaitingOrders(&models.UserContext{ID: 1, Role: enums.RoleUser})
		if err == nil {
			t.Fatalf("expected error for non-lender")
		}
	})

	t.Run("return requests error", func(t *testing.T) {
		mockReturnRequestRepo.EXPECT().
			GetAllReturnRequests([]string{return_request_status.Approved.String()}).
			Return(nil, errors.New("fail"))
		_, err := svc.GetAllApprovedAwaitingOrders(userCtx)
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("order fetch error", func(t *testing.T) {
		mockReturnRequestRepo.EXPECT().
			GetAllReturnRequests([]string{return_request_status.Approved.String()}).
			Return([]models.ReturnRequest{{OrderID: 10}}, nil)
		mockOrderRepo.EXPECT().GetOrderByID(10).Return(nil, errors.New("order fail"))
		_, err := svc.GetAllApprovedAwaitingOrders(userCtx)
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("order nil skipped", func(t *testing.T) {
		mockReturnRequestRepo.EXPECT().
			GetAllReturnRequests([]string{return_request_status.Approved.String()}).
			Return([]models.ReturnRequest{{OrderID: 10}}, nil)
		mockOrderRepo.EXPECT().GetOrderByID(10).Return(nil, nil)
		orders, err := svc.GetAllApprovedAwaitingOrders(userCtx)
		if err != nil || len(orders) != 0 {
			t.Fatalf("expected empty orders, got %v", orders)
		}
	})

	t.Run("success", func(t *testing.T) {
		mockReturnRequestRepo.EXPECT().
			GetAllReturnRequests([]string{return_request_status.Approved.String()}).
			Return([]models.ReturnRequest{{OrderID: 10}}, nil)
		mockOrderRepo.EXPECT().GetOrderByID(10).Return(&models.Order{ID: 10}, nil)
		orders, err := svc.GetAllApprovedAwaitingOrders(userCtx)
		if err != nil || len(orders) != 1 {
			t.Fatalf("expected 1 order, got %v", orders)
		}
	})
}

func TestGetLenderOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
	svc := order_service.NewOrderService(mockOrderRepo, nil, nil, logger.NewFakeLogger())

	lenderCtx := &models.UserContext{ID: 20, Role: enums.RoleLender}
	userCtx := &models.UserContext{ID: 21, Role: enums.RoleUser}

	t.Run("unauthorized user", func(t *testing.T) {
		_, err := svc.GetLenderOrders(userCtx)
		if err == nil || err.Error() != "only lender can get orders" {
			t.Fatalf("expected 'only lender can get orders', got %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			GetLenderOrders(lenderCtx.ID).
			Return(nil, errors.New("db error"))

		_, err := svc.GetLenderOrders(lenderCtx)
		if err == nil || err.Error() != "db error" {
			t.Fatalf("expected 'db error', got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.EXPECT().
			GetLenderOrders(lenderCtx.ID).
			Return([]*models.Order{{ID: 1}, {ID: 2}}, nil)

		orders, err := svc.GetLenderOrders(lenderCtx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(orders) != 2 {
			t.Fatalf("expected 2 orders, got %d", len(orders))
		}
	})
}
