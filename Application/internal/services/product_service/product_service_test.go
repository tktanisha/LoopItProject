package product_service_test

// import (
// 	"errors"
// 	"loopit/internal/enums"
// 	"loopit/internal/mock"
// 	"loopit/internal/models"
// 	"loopit/internal/services/product_service"
// 	"loopit/pkg/logger"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// )

// func TestProductService_GetAllProducts(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock.NewMockProductRepo(ctrl)
// 	log := logger.NewFakeLogger()

// 	service := product_service.NewProductService(mockRepo, nil, log)

// 	tests := []struct {
// 		name        string
// 		mockSetup   func()
// 		expectedErr bool
// 	}{
// 		{
// 			name: "success - returns products",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().FindAll().Return([]*models.ProductResponse{
// 					{Product: models.Product{ID: 1, Name: "Prod1"}},
// 				}, nil)
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "failure - repo error",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().FindAll().Return(nil, errors.New("db error"))
// 			},
// 			expectedErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()
// 			_, err := service.GetAllProducts()
// 			if (err != nil) != tt.expectedErr {
// 				t.Errorf("expectedErr=%v, got err=%v", tt.expectedErr, err)
// 			}
// 		})
// 	}
// }

// func TestProductService_GetProductByID(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock.NewMockProductRepo(ctrl)
// 	log := logger.NewFakeLogger()

// 	service := product_service.NewProductService(mockRepo, nil, log)

// 	tests := []struct {
// 		name        string
// 		id          int
// 		mockSetup   func()
// 		expectedErr bool
// 	}{
// 		{
// 			name: "invalid ID - negative",
// 			id:   -1,
// 			mockSetup: func() {
// 				// no repo call expected
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "success - found",
// 			id:   1,
// 			mockSetup: func() {
// 				mockRepo.EXPECT().FindByID(1).Return(&models.ProductResponse{
// 					Product: models.Product{ID: 1, Name: "Prod1"},
// 				}, nil)
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "failure - not found",
// 			id:   99,
// 			mockSetup: func() {
// 				mockRepo.EXPECT().FindByID(99).Return(nil, errors.New("not found"))
// 			},
// 			expectedErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()
// 			_, err := service.GetProductByID(tt.id)
// 			if (err != nil) != tt.expectedErr {
// 				t.Errorf("expectedErr=%v, got err=%v", tt.expectedErr, err)
// 			}
// 		})
// 	}
// }

// func TestProductService_CreateProduct(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock.NewMockProductRepo(ctrl)
// 	log := logger.NewFakeLogger()

// 	service := product_service.NewProductService(mockRepo, nil, log)

// 	tests := []struct {
// 		name        string
// 		product     *models.Product
// 		userCtx     *models.UserContext
// 		mockSetup   func()
// 		expectedErr bool
// 	}{
// 		{
// 			name:        "failure - user not logged in",
// 			product:     &models.Product{Name: "Prod1"},
// 			userCtx:     nil,
// 			mockSetup:   func() {},
// 			expectedErr: true,
// 		},
// 		{
// 			name:    "failure - not lender",
// 			product: &models.Product{Name: "Prod2"},
// 			userCtx: &models.UserContext{ID: 1, Role: enums.RoleUser},
// 			mockSetup: func() {
// 				// no repo call
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name:    "success - lender creates product",
// 			product: &models.Product{Name: "Prod3"},
// 			userCtx: &models.UserContext{ID: 2, Role: enums.RoleLender},
// 			mockSetup: func() {
// 				mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name:    "failure - repo create fails",
// 			product: &models.Product{Name: "Prod4"},
// 			userCtx: &models.UserContext{ID: 3, Role: enums.RoleLender},
// 			mockSetup: func() {
// 				mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
// 			},
// 			expectedErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()
// 			err := service.CreateProduct(tt.product, tt.userCtx)
// 			if (err != nil) != tt.expectedErr {
// 				t.Errorf("expectedErr=%v, got err=%v", tt.expectedErr, err)
// 			}
// 		})
// 	}
// }
