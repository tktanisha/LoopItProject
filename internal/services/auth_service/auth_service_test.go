package auth_service

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	"testing"
	"time"
)

// --- Mock Repo ---
type MockUserRepo struct {
	usersByEmail map[string]*models.User
	createCalled bool
}

func (m *MockUserRepo) FindAll() []models.User { return nil }
func (m *MockUserRepo) FindByID(userID int) (*models.User, error) {
	return nil, errors.New("not implemented")
}
func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) {
	if u, ok := m.usersByEmail[email]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}
func (m *MockUserRepo) Create(user *models.User) {
	m.createCalled = true
	m.usersByEmail[user.Email] = user
}
func (m *MockUserRepo) BecomeLender(userID int) error { return nil }
func (m *MockUserRepo) Save() error                   { return nil }

// --- Tests ---

func TestRegister_TableDriven(t *testing.T) {
	cases := []struct {
		name        string
		existing    *models.User
		newUser     *models.User
		wantErr     string
		wantCreated bool
	}{
		{
			name:        "success - new user",
			existing:    nil,
			newUser:     &models.User{Email: "new@test.com", PasswordHash: "12345"},
			wantErr:     "",
			wantCreated: true,
		},
		{
			name:        "fail - user already exists",
			existing:    &models.User{ID: 1, Email: "exists@test.com"},
			newUser:     &models.User{Email: "exists@test.com", PasswordHash: "12345"},
			wantErr:     "user already exists",
			wantCreated: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockUserRepo{usersByEmail: map[string]*models.User{}}
			if tc.existing != nil {
				mockRepo.usersByEmail[tc.existing.Email] = tc.existing
			}
			service := NewAuthService(mockRepo, logger.NewFakeLogger())

			err := service.Register(tc.newUser)

			if tc.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.wantErr != "" && (err == nil || err.Error() != tc.wantErr) {
				t.Fatalf("expected error %q, got %v", tc.wantErr, err)
			}
			if tc.wantCreated != mockRepo.createCalled {
				t.Errorf("expected createCalled=%v, got %v", tc.wantCreated, mockRepo.createCalled)
			}
			if tc.wantCreated {
				if tc.newUser.Role != enums.RoleUser {
					t.Errorf("expected default role user, got %v", tc.newUser.Role)
				}
				if tc.newUser.CreatedAt.IsZero() {
					t.Errorf("expected CreatedAt set, got zero time")
				}
			}
		})
	}
}

func TestLogin_TableDriven(t *testing.T) {
	hashed, _ := utils.HashPassword("secret")

	cases := []struct {
		name      string
		repoUser  *models.User
		email     string
		password  string
		wantErr   string
		wantToken bool
	}{
		{
			name:      "success - correct credentials",
			repoUser:  &models.User{ID: 1, Email: "user@test.com", PasswordHash: hashed, Role: enums.RoleUser, CreatedAt: time.Now()},
			email:     "user@test.com",
			password:  "secret",
			wantErr:   "",
			wantToken: true,
		},
		{
			name:      "fail - user not found",
			repoUser:  nil,
			email:     "nouser@test.com",
			password:  "secret",
			wantErr:   "invalid credentials",
			wantToken: false,
		},
		{
			name:      "fail - wrong password",
			repoUser:  &models.User{ID: 2, Email: "wrong@test.com", PasswordHash: hashed, Role: enums.RoleUser},
			email:     "wrong@test.com",
			password:  "badpass",
			wantErr:   "invalid credentials",
			wantToken: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockUserRepo{usersByEmail: map[string]*models.User{}}
			if tc.repoUser != nil {
				mockRepo.usersByEmail[tc.repoUser.Email] = tc.repoUser
			}
			service := NewAuthService(mockRepo, logger.NewFakeLogger())

			token, user, err := service.Login(tc.email, tc.password)

			if tc.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.wantErr != "" && (err == nil || err.Error() != tc.wantErr) {
				t.Fatalf("expected error %q, got %v", tc.wantErr, err)
			}
			if tc.wantToken && token == "" {
				t.Errorf("expected non-empty JWT, got empty")
			}
			if tc.wantToken && user == nil {
				t.Errorf("expected user, got nil")
			}
		})
	}
}
