package buyer_request_service_test

// import (
// 	"errors"
// 	"fmt"
// 	"loopit/internal/enums"
// 	br_status "loopit/internal/enums/buyer_request_status"
// 	"loopit/internal/mock"
// 	"loopit/internal/models"
// 	"loopit/internal/services/buyer_request_service"
// 	"loopit/pkg/logger"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// )

// func TestCreateBuyerRequest(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockBuyerRepo := mock.NewMockBuyerRequestRepo(ctrl)
// 	mockProductRepo := mock.NewMockProductRepo(ctrl)
// 	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
// 	mockCategoryRepo := mock.NewMockCategoryRepo(ctrl)

// 	svc := buyer_request_service.NewBuyerRequestService(
// 		mockBuyerRepo, mockProductRepo, mockOrderRepo, mockCategoryRepo, logger.NewFakeLogger(),
// 	)

// 	userCtx := &models.UserContext{ID: 2, Role: enums.RoleUser}

// 	tests := []struct {
// 		name          string
// 		setupMocks    func()
// 		expectedError string
// 	}{
// 		{
// 			name: "product not found",
// 			setupMocks: func() {
// 				mockProductRepo.EXPECT().FindByID(1).Return(nil, errors.New("not found"))
// 			},
// 			expectedError: "product not found",
// 		},
// 		{
// 			name: "product not available",
// 			setupMocks: func() {
// 				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
// 					Product: models.Product{IsAvailable: false},
// 				}, nil)
// 			},
// 			expectedError: "product not available",
// 		},
// 		{
// 			name: "lender requesting own product",
// 			setupMocks: func() {
// 				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
// 					Product: models.Product{IsAvailable: true, LenderID: 2},
// 				}, nil)
// 			},
// 			expectedError: "lender cannot create a buying request for their own product",
// 		},
// 		{
// 			name: "error fetching existing buyer requests",
// 			setupMocks: func() {
// 				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
// 					Product: models.Product{IsAvailable: true, LenderID: 10},
// 				}, nil)
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(gomock.Any()).Return(nil, errors.New("db fetch error"))
// 			},
// 			expectedError: "db fetch error",
// 		},
// 		{
// 			name: "duplicate request exists",
// 			setupMocks: func() {
// 				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
// 					Product: models.Product{IsAvailable: true, LenderID: 10},
// 				}, nil)
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(gomock.Any()).Return([]models.BuyingRequest{
// 					{ProductID: 1, RequestedBy: 2, Status: br_status.Pending},
// 				}, nil)
// 			},
// 			expectedError: "a pending or approved request already exists",
// 		},
// 		{
// 			name: "error while creating buyer request",
// 			setupMocks: func() {
// 				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
// 					Product: models.Product{IsAvailable: true, LenderID: 10},
// 				}, nil)
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(gomock.Any()).Return(nil, nil)
// 				mockBuyerRepo.EXPECT().CreateBuyerRequest(gomock.Any()).Return(errors.New("insert error"))
// 			},
// 			expectedError: "insert error",
// 		},
// 		{
// 			name: "success",
// 			setupMocks: func() {
// 				mockProductRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
// 					Product: models.Product{IsAvailable: true, LenderID: 10},
// 				}, nil)
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(gomock.Any()).Return(nil, nil)
// 				mockBuyerRepo.EXPECT().CreateBuyerRequest(gomock.Any()).Return(nil)
// 			},
// 			expectedError: "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMocks()
// 			err := svc.CreateBuyerRequest(1, userCtx)
// 			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
// 				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
// 			}
// 			if err != nil && tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
// 				t.Fatalf("expected %v, got %v", tt.expectedError, err)
// 			}
// 		})
// 	}
// }

// func TestUpdateBuyerRequestStatus(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockBuyerRepo := mock.NewMockBuyerRequestRepo(ctrl)
// 	mockProductRepo := mock.NewMockProductRepo(ctrl)
// 	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
// 	mockCategoryRepo := mock.NewMockCategoryRepo(ctrl)

// 	svc := buyer_request_service.NewBuyerRequestService(
// 		mockBuyerRepo, mockProductRepo, mockOrderRepo, mockCategoryRepo, logger.NewFakeLogger(),
// 	)

// 	userCtxLender := &models.UserContext{ID: 2, Role: enums.RoleLender}
// 	userCtxUser := &models.UserContext{ID: 3, Role: enums.RoleUser}

// 	tests := []struct {
// 		name          string
// 		userCtx       *models.UserContext
// 		setupMocks    func()
// 		expectedError string
// 	}{
// 		{
// 			name:    "unauthorized user",
// 			userCtx: userCtxUser,
// 			setupMocks: func() {
// 				// no mocks needed
// 			},
// 			expectedError: "unauthorized",
// 		},
// 		{
// 			name:    "invalid status",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				// no mocks needed
// 			},
// 			expectedError: "invalid status",
// 		},
// 		{
// 			name:    "buyer request not found",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return(nil, nil)
// 			},
// 			expectedError: "buyer request not found",
// 		},
// 		{
// 			name:    "reject buyer request success",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return([]models.BuyingRequest{
// 					{ID: 1, ProductID: 10, RequestedBy: 4},
// 				}, nil)
// 				mockBuyerRepo.EXPECT().UpdateStatusBuyerRequest(1, br_status.Rejected.String()).Return(nil)
// 			},
// 			expectedError: "",
// 		},
// 		{
// 			name:    "reject buyer request failure",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return([]models.BuyingRequest{
// 					{ID: 1, ProductID: 10, RequestedBy: 4},
// 				}, nil)
// 				mockBuyerRepo.EXPECT().UpdateStatusBuyerRequest(1, br_status.Rejected.String()).
// 					Return(errors.New("failed to reject"))
// 			},
// 			expectedError: "failed to reject",
// 		},
// 		{
// 			name:    "product not found while approving",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return([]models.BuyingRequest{
// 					{ID: 2, ProductID: 10, RequestedBy: 4},
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(nil, errors.New("product not found"))
// 			},
// 			expectedError: "product not found",
// 		},
// 		{
// 			name:    "category not found while approving",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return([]models.BuyingRequest{
// 					{ID: 2, ProductID: 10, RequestedBy: 4},
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Category: models.Category{ID: 5},
// 				}, nil)
// 				mockCategoryRepo.EXPECT().FindByID(5).Return(models.Category{}, errors.New("category not found"))
// 			},
// 			expectedError: "category not found",
// 		},
// 		{
// 			name:    "order creation failure",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return([]models.BuyingRequest{
// 					{ID: 2, ProductID: 10, RequestedBy: 4},
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Category: models.Category{ID: 5},
// 				}, nil)
// 				mockCategoryRepo.EXPECT().FindByID(5).Return(models.Category{
// 					Price: 100, Security: 50,
// 				}, nil)
// 				mockOrderRepo.EXPECT().CreateOrder(gomock.Any()).Return(errors.New("order creation failed"))
// 			},
// 			expectedError: "order creation failed",
// 		},
// 		{
// 			name:    "update buyer request to approved failure",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return([]models.BuyingRequest{
// 					{ID: 2, ProductID: 10, RequestedBy: 4},
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Category: models.Category{ID: 5},
// 				}, nil)
// 				mockCategoryRepo.EXPECT().FindByID(5).Return(models.Category{
// 					Price: 100, Security: 50,
// 				}, nil)
// 				mockOrderRepo.EXPECT().CreateOrder(gomock.Any()).Return(nil)
// 				mockBuyerRepo.EXPECT().UpdateStatusBuyerRequest(2, br_status.Approved.String()).
// 					Return(errors.New("failed to update status"))
// 			},
// 			expectedError: "failed to update status",
// 		},
// 		{
// 			name:    "approve buyer request success",
// 			userCtx: userCtxLender,
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests(nil).Return([]models.BuyingRequest{
// 					{ID: 2, ProductID: 10, RequestedBy: 4},
// 				}, nil)
// 				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
// 					Category: models.Category{ID: 5},
// 				}, nil)
// 				mockCategoryRepo.EXPECT().FindByID(5).Return(models.Category{
// 					Price: 100, Security: 50,
// 				}, nil)
// 				mockOrderRepo.EXPECT().CreateOrder(gomock.Any()).Return(nil)
// 				mockBuyerRepo.EXPECT().UpdateStatusBuyerRequest(2, br_status.Approved.String()).Return(nil)
// 			},
// 			expectedError: "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMocks()

// 			var status br_status.Status
// 			var requestID int

// 			switch tt.name {
// 			case "invalid status":
// 				status = br_status.Pending // invalid for update
// 				requestID = 1
// 			case "reject buyer request success", "reject buyer request failure":
// 				status = br_status.Rejected
// 				requestID = 1
// 			default:
// 				status = br_status.Approved
// 				requestID = 2
// 			}

// 			err := svc.UpdateBuyerRequestStatus(requestID, status, tt.userCtx)

// 			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
// 				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
// 			}
// 			if err != nil && tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
// 				t.Fatalf("expected %v, got %v", tt.expectedError, err)
// 			}
// 		})
// 	}
// }

// func TestGetAllBuyerRequestsByStatus(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockBuyerRepo := mock.NewMockBuyerRequestRepo(ctrl)
// 	mockProductRepo := mock.NewMockProductRepo(ctrl)
// 	mockOrderRepo := mock.NewMockOrderRepo(ctrl)
// 	mockCategoryRepo := mock.NewMockCategoryRepo(ctrl)

// 	svc := buyer_request_service.NewBuyerRequestService(
// 		mockBuyerRepo, mockProductRepo, mockOrderRepo, mockCategoryRepo, logger.NewFakeLogger(),
// 	)

// 	tests := []struct {
// 		name          string
// 		setupMocks    func()
// 		expectedLen   int
// 		expectedError string
// 	}{
// 		{
// 			name: "repo error",
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests([]string{br_status.Pending.String()}).Return(nil, errors.New("db error"))
// 			},
// 			expectedLen:   0,
// 			expectedError: "db error",
// 		},
// 		{
// 			name: "filter success",
// 			setupMocks: func() {
// 				mockBuyerRepo.EXPECT().GetAllBuyerRequests([]string{br_status.Pending.String()}).Return([]models.BuyingRequest{
// 					{ProductID: 1}, {ProductID: 2},
// 				}, nil)
// 			},
// 			expectedLen:   1,
// 			expectedError: "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMocks()
// 			res, err := svc.GetAllBuyerRequestsByStatus(1, br_status.Pending)
// 			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
// 				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
// 			}
// 			if err == nil && len(res) != tt.expectedLen {
// 				t.Fatalf("expected len %d, got %d", tt.expectedLen, len(res))
// 			}
// 		})
// 	}
// }

// // helper
// func contains(got, want string) bool {
// 	return fmt.Sprintf("%v", got) == want || (len(got) >= len(want) && got[:len(want)] == want)
// }
