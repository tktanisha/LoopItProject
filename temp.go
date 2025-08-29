package user_repo_test

import (
	"database/sql"
	"errors"
	"loopit/internal/enums"
	"loopit/internal/mock"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type fakeRow struct {
	values []interface{}
	err    error
	// read   bool
}

func (r *fakeRow) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("fakeRow.Scan: dest and values length mismatch")
	}
	for i := range dest {
		switch d := dest[i].(type) {
		case *int:
			val, ok := r.values[i].(int)
			if !ok {
				return errors.New("fakeRow.Scan: type assertion to int failed")
			}
			*d = val
		case *string:
			val, ok := r.values[i].(string)
			if !ok {
				return errors.New("fakeRow.Scan: type assertion to string failed")
			}
			*d = val
		case *time.Time:
			val, ok := r.values[i].(time.Time)
			if !ok {
				return errors.New("fakeRow.Scan: type assertion to time.Time failed")
			}
			*d = val
		default:
			return errors.New("fakeRow.Scan: unsupported scan type")
		}
	}
	return nil
}

type fakeRows struct {
	columns []string
	data    [][]interface{}
	pos     int
	err     error
}

func (r *fakeRows) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	if r.pos == 0 || r.pos > len(r.data) {
		return errors.New("fakeRows.Scan: no row available")
	}
	values := r.data[r.pos-1]
	if len(dest) != len(values) {
		return errors.New("fakeRows.Scan: length mismatch")
	}

	for i := range dest {
		switch d := dest[i].(type) {
		case *int:
			val, ok := values[i].(int)
			if !ok {
				return errors.New("fakeRows.Scan: type assertion to int failed")
			}
			*d = val
		case *string:
			val, ok := values[i].(string)
			if !ok {
				return errors.New("fakeRows.Scan: type assertion to string failed")
			}
			*d = val
		case *time.Time:
			val, ok := values[i].(time.Time)
			if !ok {
				return errors.New("fakeRows.Scan: type assertion to time.Time failed")
			}
			*d = val
		default:
			return errors.New("fakeRows.Scan: unsupported scan type")
		}
	}
	return nil
}

func fakeRowFromValues(values []interface{}) *fakeRow {
	return &fakeRow{values: values}
}

func fakeRowReturningID(id int) *fakeRow {
	return &fakeRow{values: []interface{}{id}}
}

func fakeRowsFromValues(columns []string, rows ...[]interface{}) *fakeRows {
	return &fakeRows{columns: columns, data: rows}
}

// --- Tests ---

func TestFindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockLenderRepo := mock.NewMockLenderRepo(ctrl)
	log := logger.NewFakeLogger()
	repo := user_repo.NewUserDBRepo(mockDB, mockLenderRepo, log)

	now := time.Now()
	tests := []struct {
		name     string
		mockFunc func()
		wantLen  int
	}{
		{
			name: "successfully returns multiple users",
			mockFunc: func() {
				columns := []string{"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at"}
				rows := fakeRowsFromValues(columns,
					[]interface{}{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", now},
					[]interface{}{2, "Jane", "jane@example.com", "67890", "Addr2", "hash2", 20, "lender", now},
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
				rows := fakeRowsFromValues(columns,
					[]interface{}{"invalidID", "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", now},
				)
				mockDB.EXPECT().Query(gomock.Any()).Return(rows, nil)
			},
			wantLen: 0,
		},
		{
			name: "invalid role parsing",
			mockFunc: func() {
				columns := []string{"id", "full_name", "email", "phone_number", "address", "password_hash", "society_id", "role", "created_at"}
				rows := fakeRowsFromValues(columns,
					[]interface{}{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "invalidRole", now},
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
	validRow := fakeRowFromValues([]interface{}{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", now})

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
				mockDB.EXPECT().QueryRow(gomock.Any(), 1).Return(&sql.Row{}) // will fail on Scan
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			mockFunc: func() {
				row := fakeRowFromValues([]interface{}{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "invalidRole", now})
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
	validRow := fakeRowFromValues([]interface{}{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "admin", now})

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
				row := fakeRowFromValues([]interface{}{1, "John", "john@example.com", "12345", "Addr", "hash", 10, "invalidRole", now})
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
				row := fakeRowReturningID(123)
				mockDB.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			wantPanic: false,
		},
		{
			name: "insert fails",
			mockFunc: func() {
				mockDB.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&sql.Row{})
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
