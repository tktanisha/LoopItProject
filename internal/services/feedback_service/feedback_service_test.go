package feedback_service_test

import (
	"errors"
	"fmt"
	"loopit/internal/enums/order_status"
	"loopit/internal/mock"
	"loopit/internal/models"
	"loopit/internal/services/feedback_service"

	"loopit/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGiveFeedback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := mock.NewMockFeedbackRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepo(ctrl)
	mockOrderRepo := mock.NewMockOrderRepo(ctrl)

	// Use your FakeLogger
	svc := feedback_service.NewFeedbackService(
		mockFeedbackRepo,
		mockProductRepo,
		mockOrderRepo,
		logger.NewFakeLogger(),
	)

	userCtx := &models.UserContext{ID: 2}

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
			name: "order not returned",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, ProductID: 10, Status: order_status.InUse,
				}, nil)
			},
			expectedError: "feedback can only be given for returned orders",
		},
		{
			name: "product not found",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, ProductID: 10, Status: order_status.Returned,
				}, nil)
				mockProductRepo.EXPECT().FindByID(10).Return(nil, errors.New("product missing"))
			},
			expectedError: "product missing",
		},
		{
			name: "cannot give feedback to self",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, ProductID: 10, Status: order_status.Returned,
				}, nil)
				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 2},
				}, nil)
			},
			expectedError: "you cannot give feedback to yourself",
		},
		{
			name: "repo save error",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, ProductID: 10, Status: order_status.Returned,
				}, nil)
				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 3},
				}, nil)
				mockFeedbackRepo.EXPECT().CreateFeedback(gomock.Any()).Return(errors.New("db error"))
			},
			expectedError: "db error",
		},
		{
			name: "success",
			setupMocks: func() {
				mockOrderRepo.EXPECT().GetOrderByID(1).Return(&models.Order{
					ID: 1, ProductID: 10, Status: order_status.Returned,
				}, nil)
				mockProductRepo.EXPECT().FindByID(10).Return(&models.ProductResponse{
					Product: models.Product{LenderID: 3},
				}, nil)
				mockFeedbackRepo.EXPECT().CreateFeedback(gomock.Any()).DoAndReturn(
					func(f models.Feedback) error {
						if f.GivenBy != 2 || f.GivenTo != 3 || f.Text != "Nice product" {
							return fmt.Errorf("unexpected feedback %+v", f)
						}
						if f.CreatedAt.IsZero() {
							return fmt.Errorf("CreatedAt not set")
						}
						return nil
					},
				)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := svc.GiveFeedback(1, "Nice product", 5, userCtx)
			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
				t.Fatalf("expected error %q, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetAllGivenFeedbacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := mock.NewMockFeedbackRepository(ctrl)

	svc := feedback_service.NewFeedbackService(
		mockFeedbackRepo, nil, nil, logger.NewFakeLogger(),
	)

	userCtx := &models.UserContext{ID: 2}

	tests := []struct {
		name          string
		setupMocks    func()
		expectedCount int
		expectErr     bool
	}{
		{
			name: "repo error",
			setupMocks: func() {
				mockFeedbackRepo.EXPECT().GetAllFeedbacks().Return(nil, errors.New("db fail"))
			},
			expectedCount: 0,
			expectErr:     true,
		},
		{
			name: "filter given feedbacks",
			setupMocks: func() {
				mockFeedbackRepo.EXPECT().GetAllFeedbacks().Return([]models.Feedback{
					{ID: 1, GivenBy: 2},
					{ID: 2, GivenBy: 3},
				}, nil)
			},
			expectedCount: 1,
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			feedbacks, err := svc.GetAllGivenFeedbacks(userCtx)
			if tt.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectErr && len(feedbacks) != tt.expectedCount {
				t.Fatalf("expected %d feedbacks, got %d", tt.expectedCount, len(feedbacks))
			}
		})
	}
}

func TestGetAllReceivedFeedbacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedbackRepo := mock.NewMockFeedbackRepository(ctrl)

	svc := feedback_service.NewFeedbackService(
		mockFeedbackRepo, nil, nil, logger.NewFakeLogger(),
	)

	userCtx := &models.UserContext{ID: 5}

	tests := []struct {
		name          string
		setupMocks    func()
		expectedCount int
		expectErr     bool
	}{
		{
			name: "repo error",
			setupMocks: func() {
				mockFeedbackRepo.EXPECT().GetAllFeedbacks().Return(nil, errors.New("db fail"))
			},
			expectedCount: 0,
			expectErr:     true,
		},
		{
			name: "filter received feedbacks",
			setupMocks: func() {
				mockFeedbackRepo.EXPECT().GetAllFeedbacks().Return([]models.Feedback{
					{ID: 1, GivenTo: 5},
					{ID: 2, GivenTo: 3},
				}, nil)
			},
			expectedCount: 1,
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			feedbacks, err := svc.GetAllReceivedFeedbacks(userCtx)
			if tt.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectErr && len(feedbacks) != tt.expectedCount {
				t.Fatalf("expected %d feedbacks, got %d", tt.expectedCount, len(feedbacks))
			}
		})
	}
}
