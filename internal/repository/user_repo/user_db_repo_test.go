package user_repo_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"loopit/internal/enums"
	"loopit/internal/mock"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestFindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockLenderRepo := mock.NewMockLenderRepo(ctrl)
	log := logger.NewFakeLogger()

	repo := user_repo.NewUserDBRepo(mockDB, mockLenderRepo, log)

	tests := []struct {
		name     string
		mockFunc func()
		wantLen  int
	}{
		{
			name: "successfully returns multiple users",
			mockFunc: func() {

				columns := []string{"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at"}
				rows := sqlmockRows(columns,
					[]driver.Value{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", time.Now()},
					[]driver.Value{2, "Jane", "jane@example.com", "67890", "Addr2", "hash2", 20, "lender", time.Now()},
				)
				mockDB.EXPECT().Query(gomock.Any()).Return(rows, nil)
			},
			wantLen: 2,
		},
		{
			name: "query fails",
			mockFunc: func() {
				mockDB.EXPECT().Query(gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantLen: 0,
		},
		{
			name: "row scan fails",
			mockFunc: func() {
				columns := []string{"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at"}
				// mismatch types cause scan failure
				rows := sqlmockRows(columns,
					[]driver.Value{"invalidID", "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", time.Now()},
				)
				mockDB.EXPECT().Query(gomock.Any()).Return(rows, nil)
			},
			wantLen: 0,
		},
		{
			name: "invalid role parsing",
			mockFunc: func() {
				columns := []string{"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at"}
				rows := sqlmockRows(columns,
					[]driver.Value{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "invalidRole", time.Now()},
				)
				mockDB.EXPECT().Query(gomock.Any()).Return(rows, nil)
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			users := repo.FindAll()
			if len(users) != tt.wantLen {
				t.Errorf("expected %d users, got %d", tt.wantLen, len(users))
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockLenderRepo := mock.NewMockLenderRepo(ctrl)
	log := logger.NewFakeLogger()
	repo := user_repo.NewUserDBRepo(mockDB, mockLenderRepo, log)

	now := time.Now()
	validRow := sqlmockRow([]driver.Value{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", now})

	tests := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "user found",
			mockFunc: func() {
				mockDB.EXPECT().QueryRow(gomock.Any(), 1).Return(validRow)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			mockFunc: func() {
				row := &sql.Row{}
				mockDB.EXPECT().QueryRow(gomock.Any(), 1).Return(row)
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			mockFunc: func() {
				row := sqlmockRow([]driver.Value{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "invalidRole", now})
				mockDB.EXPECT().QueryRow(gomock.Any(), 1).Return(row)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			_, err := repo.FindByID(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestFindByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockLenderRepo := mock.NewMockLenderRepo(ctrl)
	log := logger.NewFakeLogger()
	repo := user_repo.NewUserDBRepo(mockDB, mockLenderRepo, log)

	now := time.Now()
	validRow := sqlmockRow([]driver.Value{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", now})

	tests := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "user found",
			mockFunc: func() {
				mockDB.EXPECT().QueryRow(gomock.Any(), "john@example.com").Return(validRow)
			},
			wantErr: false,
		},
		{
			name: "no user found",
			mockFunc: func() {
				mockDB.EXPECT().QueryRow(gomock.Any(), "john@example.com").Return(&sql.Row{})
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			mockFunc: func() {
				row := sqlmockRow([]driver.Value{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "invalidRole", now})
				mockDB.EXPECT().QueryRow(gomock.Any(), "john@example.com").Return(row)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			_, err := repo.FindByEmail("john@example.com")
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockLenderRepo := mock.NewMockLenderRepo(ctrl)
	log := logger.NewFakeLogger()
	repo := user_repo.NewUserDBRepo(mockDB, mockLenderRepo, log)

	user := &models.User{FullName: "John", Email: "john@example.com", Role: enums.RoleAdmin}

	tests := []struct {
		name      string
		mockFunc  func()
		wantPanic bool
	}{
		{
			name: "success",
			mockFunc: func() {
				row := sqlmockRowReturningID(123)
				mockDB.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			wantPanic: false,
		},
		{
			name: "insert fails",
			mockFunc: func() {
				mockDB.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&sql.Row{})
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("expected panic=%v, got %v", tt.wantPanic, r)
				}
			}()
			repo.Create(user)
		})
	}
}

func TestBecomeLender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockLenderRepo := mock.NewMockLenderRepo(ctrl)
	log := logger.NewFakeLogger()
	repo := user_repo.NewUserDBRepo(mockDB, mockLenderRepo, log)

	tests := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				mockDB.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockLenderRepo.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "db update fails",
			mockFunc: func() {
				mockDB.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "lender creation fails",
			mockFunc: func() {
				mockDB.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockLenderRepo.EXPECT().Create(gomock.Any()).Return(errors.New("insert error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.BecomeLender(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockLenderRepo := mock.NewMockLenderRepo(ctrl)
	log := logger.NewFakeLogger()
	repo := user_repo.NewUserDBRepo(mockDB, mockLenderRepo, log)

	if err := repo.Save(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

// --- Helper Functions ---

// sqlmockRows creates a *sql.Rows for testing
func sqlmockRows(columns []string, rows ...[]driver.Value) *sqlmock.Rows {
	rowSet := sqlmock.NewRows(columns)
	for _, r := range rows {
		rowSet.AddRow(r...)
	}
	return rowSet
}

// sqlmockRow creates a *sql.Row from values
func sqlmockRow(values []driver.Value) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at",
	}).AddRow(values...)
	return rows
}

// sqlmockRowReturningID creates a *sql.Row returning only ID for insert
func sqlmockRowReturningID(id int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
	return rows
}
