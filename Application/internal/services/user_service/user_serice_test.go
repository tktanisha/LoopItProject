package user_service

// import (
// 	"errors"
// 	"loopit/internal/enums"
// 	"loopit/internal/models"
// 	"loopit/pkg/logger"
// 	"testing"
// )

// // MockUserRepo used for testing
// type MockUserRepo struct {
// 	BecomeLenderCalledWith int
// 	BecomeLenderErr        error
// }

// func (m *MockUserRepo) FindAll() []models.User                         { return nil }
// func (m *MockUserRepo) FindByID(userID int) (*models.User, error)      { return nil, nil }
// func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) { return nil, nil }
// func (m *MockUserRepo) Create(user *models.User)                       {}
// func (m *MockUserRepo) BecomeLender(userID int) error {
// 	m.BecomeLenderCalledWith = userID
// 	return m.BecomeLenderErr
// }
// func (m *MockUserRepo) Save() error { return nil }

// func TestBecomeLender_TableDriven(t *testing.T) {
// 	cases := []struct {
// 		name         string
// 		user         *models.UserContext
// 		repoErr      error
// 		wantErr      string
// 		wantCallRepo bool
// 	}{
// 		{
// 			name:         "already lender",
// 			user:         &models.UserContext{ID: 1, Role: enums.RoleLender},
// 			repoErr:      nil,
// 			wantErr:      "user is already a lender",
// 			wantCallRepo: false,
// 		},
// 		{
// 			name:         "become lender success",
// 			user:         &models.UserContext{ID: 2, Role: enums.RoleUser},
// 			repoErr:      nil,
// 			wantErr:      "",
// 			wantCallRepo: true,
// 		},
// 		{
// 			name:         "become lender repo failure",
// 			user:         &models.UserContext{ID: 3, Role: enums.RoleUser},
// 			repoErr:      errors.New("repo fail"),
// 			wantErr:      "repo fail",
// 			wantCallRepo: true,
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			repo := &MockUserRepo{BecomeLenderErr: tc.repoErr}
// 			fakeLogger := logger.NewFakeLogger()
// 			service := NewUserService(repo, fakeLogger)

// 			err := service.BecomeLender(tc.user)

// 			if tc.wantErr == "" && err != nil {
// 				t.Fatalf("expected no error, got %v", err)
// 			}
// 			if tc.wantErr != "" {
// 				if err == nil || err.Error() != tc.wantErr {
// 					t.Fatalf("expected error %q, got %v", tc.wantErr, err)
// 				}
// 			}
// 			if tc.wantCallRepo {
// 				if repo.BecomeLenderCalledWith != tc.user.ID {
// 					t.Errorf("expected repo.BecomeLender called with %d, got %d", tc.user.ID, repo.BecomeLenderCalledWith)
// 				}
// 			} else {
// 				if repo.BecomeLenderCalledWith != 0 {
// 					t.Errorf("expected repo.BecomeLender NOT called but got call with %d", repo.BecomeLenderCalledWith)
// 				}
// 			}
// 		})
// 	}
// }
