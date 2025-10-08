package user_repo_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"loopit/internal/db"
	"loopit/internal/enums"
	custom_mock "loopit/internal/mock"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestUserDBRepo_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		rows      *sqlmock.Rows
		expectErr bool
		wantCount int
	}{
		{
			name: "successfully returns multiple users",
			rows: sqlmock.NewRows([]string{
				"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
			}).AddRow(1, "John", "john@example.com", "12345", "NY", "hash", 10, enums.RoleLender.String(), time.Now()),
			expectErr: false,
			wantCount: 1,
		},
		{
			name:      "query fails",
			rows:      nil,
			expectErr: true,
			wantCount: 0,
		},
		{
			name: "invalid role parsing skips user",
			rows: sqlmock.NewRows([]string{
				"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
			}).AddRow(2, "Jane", "jane@example.com", "54321", "LA", "hash", 11, "invalid_role", time.Now()),
			expectErr: false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			query := "SELECT id, full_name, email, phone_number, address, password_hash, society_id, role, created_at FROM users"

			if tt.rows != nil {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(tt.rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(sql.ErrConnDone)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockLender := custom_mock.NewMockLenderRepo(ctrl)
			logger := logger.NewFakeLogger()
			repo := user_repo.NewUserDBRepo(db.NewRealDatabase(sqlDB), mockLender, logger)

			users := repo.FindAll()

			if len(users) != tt.wantCount {
				t.Errorf("expected %d users, got %d", tt.wantCount, len(users))
			}
		})
	}
}

func TestUserDBRepo_FindByID(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	query := "SELECT id, full_name, email, phone_number, address, password_hash, society_id, role, created_at FROM users WHERE id=$1"

	tests := []struct {
		name      string
		userID    int
		rows      *sqlmock.Rows
		expectErr bool
	}{
		{
			name:   "success",
			userID: 1,
			rows: sqlmock.NewRows([]string{
				"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
			}).AddRow(1, "John", "john@example.com", "12345", "NY", "hash", 10, enums.RoleLender.String(), time.Now()),
			expectErr: false,
		},
		{
			name:      "not found",
			userID:    2,
			rows:      sqlmock.NewRows([]string{"id", "full_name"}), // empty
			expectErr: true,
		},
		{
			name:   "scan error",
			userID: 3,
			rows: sqlmock.NewRows([]string{
				"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
			}).AddRow("bad_id", "John", "john@example.com", "12345", "NY", "hash", 10, enums.RoleLender.String(), time.Now()),
			expectErr: true,
		},
		{
			name:   "invalid role parsing",
			userID: 4,
			rows: sqlmock.NewRows([]string{
				"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
			}).AddRow(4, "Jane", "jane@example.com", "54321", "LA", "hash", 11, "invalid_role", time.Now()),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tt.userID).WillReturnRows(tt.rows)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockLender := custom_mock.NewMockLenderRepo(ctrl)
			logger := logger.NewFakeLogger()
			repo := user_repo.NewUserDBRepo(db.NewRealDatabase(sqlDB), mockLender, logger)

			user, err := repo.FindByID(tt.userID)

			if tt.expectErr && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("did not expect error, got: %v", err)
			}
			if !tt.expectErr && user == nil {
				t.Errorf("expected user, got nil")
			}
		})
	}
}

func TestUserDBRepo_FindByEmail(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	query := "SELECT id, full_name, email, phone_number, address, password_hash, society_id, role, created_at FROM users WHERE email=$1"

	tests := []struct {
		name    string
		email   string
		rows    *sqlmock.Rows
		wantErr bool
		wantNil bool
	}{
		{
			name:  "success",
			email: "john@example.com",
			rows: sqlmock.NewRows([]string{
				"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
			}).AddRow(1, "John", "john@example.com", "12345", "NY", "hash", 10, enums.RoleLender.String(), time.Now()),
			wantErr: false,
			wantNil: false,
		},
		{
			name:    "not found",
			email:   "notfound@example.com",
			rows:    sqlmock.NewRows([]string{"id"}), // empty
			wantErr: true,
			wantNil: true,
		},
		{
			name:  "db scan error",
			email: "broken@example.com",
			rows: sqlmock.NewRows([]string{
				"id", "full_name", // missing required columns
			}).AddRow(1, "John"),
			wantErr: true,
			wantNil: true,
		},
		{
			name:  "invalid role parsing",
			email: "invalidrole@example.com",
			rows: sqlmock.NewRows([]string{
				"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
			}).AddRow(1, "John", "invalidrole@example.com", "12345", "NY", "hash", 10, "INVALID_ROLE", time.Now()),
			wantErr: true,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tt.email).WillReturnRows(tt.rows)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockLender := custom_mock.NewMockLenderRepo(ctrl)
			logger := logger.NewFakeLogger()
			repo := user_repo.NewUserDBRepo(db.NewRealDatabase(sqlDB), mockLender, logger)

			user, err := repo.FindByEmail(tt.email)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantNil && user != nil {
				t.Errorf("expected nil user, got %+v", user)
			}
		})
	}
}

func TestUserDBRepo_Create(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	query := `
	INSERT INTO users (full_name, email, phone_number, address, password_hash, society_id, role, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`

	tests := []struct {
		name        string
		returnError bool
	}{
		{name: "success", returnError: false},
		{name: "db insert error", returnError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.returnError {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("John", "john@example.com", "12345", "NY", "hash", 10, enums.RoleUser.String(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("John", "john@example.com", "12345", "NY", "hash", 10, enums.RoleUser.String(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockLender := custom_mock.NewMockLenderRepo(ctrl)
			logger := logger.NewFakeLogger()
			repo := user_repo.NewUserDBRepo(db.NewRealDatabase(sqlDB), mockLender, logger)

			user := &models.User{
				FullName:     "John",
				Email:        "john@example.com",
				PhoneNumber:  "12345",
				Address:      "NY",
				PasswordHash: "hash",
				SocietyID:    10,
				Role:         enums.RoleUser,
			}

			if tt.returnError {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic, got none")
					}
				}()
			}

			repo.Create(user)

			if !tt.returnError && user.ID != 1 {
				t.Errorf("expected user ID 1, got %d", user.ID)
			}
		})
	}
}

func TestUserDBRepo_BecomeLender(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	tests := []struct {
		name          string
		updateError   bool
		lenderError   bool
		expectedError bool
	}{
		{name: "success", updateError: false, lenderError: false, expectedError: false},
		{name: "db update failure", updateError: true, lenderError: false, expectedError: true},
		{name: "lender create failure", updateError: false, lenderError: true, expectedError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.updateError {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET role=$1 WHERE id=$2")).
					WithArgs(enums.RoleLender.String(), 1).
					WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET role=$1 WHERE id=$2")).
					WithArgs(enums.RoleLender.String(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLender := custom_mock.NewMockLenderRepo(ctrl)
			if !tt.updateError {
				if tt.lenderError {
					mockLender.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("failed to create lender")).Times(1)
				} else {
					mockLender.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
				}
			}

			logger := logger.NewFakeLogger()
			repo := user_repo.NewUserDBRepo(db.NewRealDatabase(sqlDB), mockLender, logger)

			err := repo.BecomeLender(1)
			if tt.expectedError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
