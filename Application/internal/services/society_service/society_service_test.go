// file: internal/services/society_service/society_service_test.go
package society_service_test

// import (
// 	"errors"
// 	"loopit/internal/mock"
// 	"loopit/internal/models"
// 	"loopit/internal/services/society_service"
// 	"loopit/pkg/logger"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// )

// func TestGetAllSocieties(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	log := logger.NewFakeLogger()
// 	mockRepo := mock.NewMockSocietyRepo(ctrl)
// 	svc := society_service.NewSocietyService(mockRepo, log)

// 	tests := []struct {
// 		name       string
// 		mockSetup  func()
// 		wantErr    bool
// 		wantLength int
// 	}{
// 		{
// 			name: "success",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().
// 					FindAll().
// 					Return([]models.Society{{ID: 1, Name: "GreenVille"}}, nil)
// 			},
// 			wantErr:    false,
// 			wantLength: 1,
// 		},
// 		{
// 			name: "repo error",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().
// 					FindAll().
// 					Return(nil, errors.New("db failed"))
// 			},
// 			wantErr:    true,
// 			wantLength: 0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()
// 			societies, err := svc.GetAllSocieties()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("expected err=%v, got %v", tt.wantErr, err)
// 			}
// 			if len(societies) != tt.wantLength {
// 				t.Errorf("expected len=%d, got %d", tt.wantLength, len(societies))
// 			}
// 		})
// 	}
// }

// func TestCreateSociety(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	log := logger.NewFakeLogger()
// 	mockRepo := mock.NewMockSocietyRepo(ctrl)
// 	svc := society_service.NewSocietyService(mockRepo, log)

// 	tests := []struct {
// 		name      string
// 		mockSetup func()
// 		wantErr   bool
// 	}{
// 		{
// 			name: "success",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().
// 					Create(gomock.Any()).
// 					Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repo error",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().
// 					Create(gomock.Any()).
// 					Return(errors.New("insert failed"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()
// 			err := svc.CreateSociety("GreenVille", "CityCenter", "12345")
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("expected err=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }
