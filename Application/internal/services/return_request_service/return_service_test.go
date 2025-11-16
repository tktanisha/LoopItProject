package return_request_service_test

// import (
// 	"errors"
// 	"testing"

// 	"loopit/internal/enums/order_status"
// 	"loopit/internal/enums/return_request_status"
// 	"loopit/internal/mock"
// 	"loopit/internal/models"
// 	"loopit/internal/services/return_request_service"

// 	"loopit/pkg/logger"

// 	"github.com/golang/mock/gomock"
// )

// func TestCreateReturnRequest(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
// 	mockProductRepo := mock.NewMockProductRepo(ctrl)
// 	mockRRRepo := mock.NewMockReturnRequestRepo(ctrl)
// 	log := logger.NewFakeLogger()

// 	service := return_request_service.NewReturnRequestService(mockOrderRepo, mockProductRepo, mockRRRepo, log)

// 	tests := []struct {
// 		name      string
// 		setup     func()
// 		expectErr bool
// 	}{
// 		{
// 			name: "success - valid return request",
// 			setup: func() {
// 				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
// 					ID:        1,
// 					ProductID: 10,
// 					UserID:    2,
// 					Status:    order_status.InUse,
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Product: models.Product{ID: 10, LenderID: 2},
// 				}, nil)
// 				mockOrderRepo.EXPECT().UpdateOrderStatus(1, order_status.ReturnRequested.String()).Return(nil)
// 				mockRRRepo.EXPECT().CreateReturnRequest(gomock.Any()).Return(nil)
// 			},
// 			expectErr: false,
// 		},
// 		{
// 			name: "fail - order not in use",
// 			setup: func() {
// 				mockOrderRepo.EXPECT().GetOrderByID(1).
// 					Return(&models.Order{ID: 1, Status: order_status.Returned}, nil)
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name: "fail - repo error fetching order",
// 			setup: func() {
// 				mockOrderRepo.EXPECT().GetOrderByID(1).Return(nil, errors.New("db error"))
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name: "fail - product repo error",
// 			setup: func() {
// 				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
// 					ID:        1,
// 					ProductID: 10,
// 					Status:    order_status.InUse,
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(nil, errors.New("product fetch fail"))
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name: "fail - user not lender",
// 			setup: func() {
// 				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
// 					ID:        1,
// 					ProductID: 10,
// 					Status:    order_status.InUse,
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Product: models.Product{ID: 10, LenderID: 999}, // mismatch
// 				}, nil)
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name: "fail - update order status error",
// 			setup: func() {
// 				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
// 					ID:        1,
// 					ProductID: 10,
// 					Status:    order_status.InUse,
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Product: models.Product{ID: 10, LenderID: 2},
// 				}, nil)
// 				mockOrderRepo.EXPECT().UpdateOrderStatus(1, order_status.ReturnRequested.String()).
// 					Return(errors.New("update fail"))
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name: "fail - create return request error",
// 			setup: func() {
// 				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
// 					ID:        1,
// 					ProductID: 10,
// 					Status:    order_status.InUse,
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Product: models.Product{ID: 10, LenderID: 2},
// 				}, nil)
// 				mockOrderRepo.EXPECT().UpdateOrderStatus(1, order_status.ReturnRequested.String()).Return(nil)
// 				mockRRRepo.EXPECT().CreateReturnRequest(gomock.Any()).
// 					Return(errors.New("create fail"))
// 			},
// 			expectErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
// 			err := service.CreateReturnRequest(2, 1)
// 			if tt.expectErr && err == nil {
// 				t.Errorf("expected error, got nil")
// 			}
// 			if !tt.expectErr && err != nil {
// 				t.Errorf("expected no error, got %v", err)
// 			}
// 		})
// 	}
// }

// func TestUpdateReturnRequestStatus(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
// 	mockProductRepo := mock.NewMockProductRepo(ctrl) // Not used here but part of service
// 	mockRRRepo := mock.NewMockReturnRequestRepo(ctrl)
// 	log := logger.NewFakeLogger()

// 	service := return_request_service.NewReturnRequestService(mockOrderRepo, mockProductRepo, mockRRRepo, log)

// 	tests := []struct {
// 		name      string
// 		newStatus return_request_status.Status
// 		setup     func()
// 		expectErr bool
// 	}{
// 		{
// 			name:      "success - approve request",
// 			newStatus: return_request_status.Approved,
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetReturnRequestByID(1).Return(models.ReturnRequest{
// 					ID: 1, OrderID: 5, Status: return_request_status.Pending,
// 				}, nil)
// 				mockOrderRepo.EXPECT().GetOrderByID(5).Return(&models.Order{ID: 5, UserID: 2}, nil)
// 				mockRRRepo.EXPECT().UpdateReturnRequestStatus(1, return_request_status.Approved.String()).Return(nil)
// 			},
// 			expectErr: false,
// 		},
// 		{
// 			name:      "fail - invalid status",
// 			newStatus: return_request_status.Pending,
// 			setup:     func() {},
// 			expectErr: true,
// 		},
// 		{
// 			name:      "fail - request not pending",
// 			newStatus: return_request_status.Rejected,
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetReturnRequestByID(1).Return(models.ReturnRequest{
// 					ID: 1, OrderID: 5, Status: return_request_status.Approved,
// 				}, nil)
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name:      "fail - user does not own order",
// 			newStatus: return_request_status.Approved,
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetReturnRequestByID(1).Return(models.ReturnRequest{
// 					ID: 1, OrderID: 5, Status: return_request_status.Pending,
// 				}, nil)
// 				mockOrderRepo.EXPECT().GetOrderByID(5).Return(&models.Order{ID: 5, UserID: 99}, nil)
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name:      "fail - error fetching return request",
// 			newStatus: return_request_status.Approved,
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetReturnRequestByID(1).Return(models.ReturnRequest{}, errors.New("db error"))
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name:      "fail - error fetching order",
// 			newStatus: return_request_status.Approved,
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetReturnRequestByID(1).Return(models.ReturnRequest{
// 					ID: 1, OrderID: 5, Status: return_request_status.Pending,
// 				}, nil)
// 				mockOrderRepo.EXPECT().GetOrderByID(5).Return(nil, errors.New("order not found"))
// 			},
// 			expectErr: true,
// 		},
// 		{
// 			name:      "fail - error updating return request status",
// 			newStatus: return_request_status.Approved,
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetReturnRequestByID(1).Return(models.ReturnRequest{
// 					ID: 1, OrderID: 5, Status: return_request_status.Pending,
// 				}, nil)
// 				mockOrderRepo.EXPECT().GetOrderByID(5).Return(&models.Order{ID: 5, UserID: 2}, nil)
// 				mockRRRepo.EXPECT().UpdateReturnRequestStatus(1, return_request_status.Approved.String()).
// 					Return(errors.New("update failed"))
// 			},
// 			expectErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
// 			err := service.UpdateReturnRequestStatus(2, 1, tt.newStatus)
// 			if tt.expectErr && err == nil {
// 				t.Errorf("expected error, got nil")
// 			}
// 			if !tt.expectErr && err != nil {
// 				t.Errorf("expected no error, got %v", err)
// 			}
// 		})
// 	}
// }

// func TestGetPendingReturnRequests(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
// 	mockProductRepo := mock.NewMockProductRepo(ctrl)
// 	mockRRRepo := mock.NewMockReturnRequestRepo(ctrl)
// 	log := logger.NewFakeLogger()

// 	service := return_request_service.NewReturnRequestService(mockOrderRepo, mockProductRepo, mockRRRepo, log)

// 	tests := []struct {
// 		name      string
// 		setup     func()
// 		expectLen int
// 		expectErr bool
// 	}{
// 		{
// 			name: "success - returns only requests owned by user",
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetAllReturnRequests([]string{return_request_status.Pending.String()}).Return([]models.ReturnRequest{
// 					{ID: 1, OrderID: 5, Status: return_request_status.Pending},
// 					{ID: 2, OrderID: 6, Status: return_request_status.Pending},
// 				}, nil)

// 				mockOrderRepo.EXPECT().GetOrderByID(5).Return(&models.Order{ID: 5, UserID: 2}, nil)
// 				mockOrderRepo.EXPECT().GetOrderByID(6).Return(&models.Order{ID: 6, UserID: 99}, nil)
// 			},
// 			expectLen: 1,
// 			expectErr: false,
// 		},
// 		{
// 			name: "fail - repo error fetching return requests",
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetAllReturnRequests([]string{return_request_status.Pending.String()}).Return(nil, errors.New("db error"))
// 			},
// 			expectLen: 0,
// 			expectErr: true,
// 		},
// 		{
// 			name: "skip - order fetch error",
// 			setup: func() {
// 				mockRRRepo.EXPECT().GetAllReturnRequests([]string{return_request_status.Pending.String()}).Return([]models.ReturnRequest{
// 					{ID: 3, OrderID: 7, Status: return_request_status.Pending},
// 				}, nil)
// 				mockOrderRepo.EXPECT().GetOrderByID(7).Return(nil, errors.New("not found"))
// 			},
// 			expectLen: 0,
// 			expectErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
// 			requests, err := service.GetPendingReturnRequests(2)

// 			if tt.expectErr && err == nil {
// 				t.Errorf("expected error, got nil")
// 			}
// 			if !tt.expectErr && err != nil {
// 				t.Errorf("expected no error, got %v", err)
// 			}
// 			if !tt.expectErr && len(requests) != tt.expectLen {
// 				t.Errorf("expected %d requests, got %d", tt.expectLen, len(requests))
// 			}
// 		})
// 	}
// }
