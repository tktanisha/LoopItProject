package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"loopit/internal/db"
	"loopit/internal/mock"

	"github.com/golang/mock/gomock"
)

// Test Query
func TestPostgresDB_Query(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockRows := &sql.Rows{}

	mockDB.EXPECT().
		Query("SELECT * FROM users WHERE id=$1", 1).
		Return(mockRows, nil)

	pg := db.NewPostgresDB(mockDB)

	_, err := pg.Query("SELECT * FROM users WHERE id=$1", 1)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
}

// Test QueryRow
func TestPostgresDB_QueryRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)
	mockRow := &sql.Row{}

	mockDB.EXPECT().
		QueryRow("SELECT name FROM users WHERE id=$1", 1).
		Return(mockRow)

	pg := db.NewPostgresDB(mockDB)

	row := pg.QueryRow("SELECT name FROM users WHERE id=$1", 1)
	if row == nil {
		t.Fatalf("Expected non-nil row")
	}
}

// Test Exec
type fakeResult struct{}

func (f fakeResult) LastInsertId() (int64, error) { return 1, nil }
func (f fakeResult) RowsAffected() (int64, error) { return 1, nil }

func TestPostgresDB_Exec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)

	mockDB.EXPECT().
		Exec("DELETE FROM users WHERE id=$1", 1).
		Return(fakeResult{}, nil)

	pg := db.NewPostgresDB(mockDB)

	_, err := pg.Exec("DELETE FROM users WHERE id=$1", 1)
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
}

// Test Close
func TestPostgresDB_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabaseInterface(ctrl)

	mockDB.EXPECT().
		Close().
		Return(nil)

	pg := db.NewPostgresDB(mockDB)

	if err := pg.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestConnectDB(t *testing.T) {
	tests := []struct {
		name      string
		connStr   string
		mockOpen  func() func(driverName, dataSourceName string) (*sql.DB, error)
		mockPing  func() func(db *sql.DB) error
		expectErr bool
	}{
		{
			name:    "Connection success",
			connStr: "postgresql://neondb_owner:npg_V8PhgoDLN5dq@ep-lingering-mode-ae1htbuw-pooler.c-2.us-east-2.aws.neon.tech/neondb?sslmode=require&channel_binding=require",
			mockOpen: func() func(driverName, dataSourceName string) (*sql.DB, error) {
				return func(_, _ string) (*sql.DB, error) {
					return &sql.DB{}, nil
				}
			},
			mockPing: func() func(db *sql.DB) error {
				return func(_ *sql.DB) error { return nil }
			},
			expectErr: false,
		},
		{
			name:    "Open failure",
			connStr: "invalid",
			mockOpen: func() func(driverName, dataSourceName string) (*sql.DB, error) {
				return func(_, _ string) (*sql.DB, error) {
					return nil, errors.New("open failed")
				}
			},
			mockPing: func() func(db *sql.DB) error {
				return func(_ *sql.DB) error { return nil }
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origOpen := db.SqlOpen
			origPing := db.SqlPing
			defer func() {
				db.SqlOpen = origOpen
				db.SqlPing = origPing
			}()

			db.SqlOpen = tt.mockOpen()
			db.SqlPing = tt.mockPing()

			_, err := db.ConnectDB(tt.connStr)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}
