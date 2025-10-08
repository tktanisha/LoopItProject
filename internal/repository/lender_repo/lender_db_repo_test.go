package lender_repo_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/lender_repo"
	"loopit/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLenderDBRepo_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		wantCount int
	}{
		{
			name: "success - multiple lenders",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "is_verified", "total_earnings"}).
					AddRow(1, true, 1000).
					AddRow(2, false, 200)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, is_verified, total_earnings FROM lenders")).
					WillReturnRows(rows)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, is_verified, total_earnings FROM lenders")).
					WillReturnError(errors.New("db failure"))
			},
			wantErr:   true,
			wantCount: 0,
		},
		{
			name: "scan error - row skipped",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "is_verified", "total_earnings"}).
					AddRow(nil, true, 100) // invalid ID to force scan error
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, is_verified, total_earnings FROM lenders")).
					WillReturnRows(rows)
			},
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()
			tt.mockSetup(mock)

			repo := lender_repo.NewLenderDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
			lenders, err := repo.FindAll()

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() expected error=%v, got %v", tt.wantErr, err)
			}

			if len(lenders) != tt.wantCount {
				t.Errorf("FindAll() expected count=%d, got %d", tt.wantCount, len(lenders))
			}
		})
	}
}

func TestLenderDBRepo_FindByID(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		wantNil   bool
	}{
		{
			name:   "success - lender found",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "is_verified", "total_earnings"}).
					AddRow(1, true, 500)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, is_verified, total_earnings FROM lenders WHERE id=$1")).
					WithArgs(1).WillReturnRows(rows)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name:   "lender not found",
			userID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, is_verified, total_earnings FROM lenders WHERE id=$1")).
					WithArgs(2).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name:   "db error",
			userID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, is_verified, total_earnings FROM lenders WHERE id=$1")).
					WithArgs(3).WillReturnError(errors.New("db failure"))
			},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()
			tt.mockSetup(mock)

			repo := lender_repo.NewLenderDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
			lender, err := repo.FindByID(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() expected error=%v, got %v", tt.wantErr, err)
			}

			if (lender == nil) != tt.wantNil {
				t.Errorf("FindByID() expected nil=%v, got %v", tt.wantNil, lender)
			}
		})
	}
}

func TestLenderDBRepo_Create(t *testing.T) {
	tests := []struct {
		name      string
		lender    *models.Lender
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:   "success - lender created",
			lender: &models.Lender{ID: 1, IsVerified: true, TotalEarnings: 100.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM lenders WHERE id=$1)")).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
				mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO lenders (id, is_verified, total_earnings)
    VALUES ($1, $2, $3)
    RETURNING id
    `)).
					WithArgs(1, true, 100.0).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name:   "already exists - skip insert",
			lender: &models.Lender{ID: 2, IsVerified: false, TotalEarnings: 200.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM lenders WHERE id=$1)")).
					WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			wantErr: false,
		},
		{
			name:   "existence check fails",
			lender: &models.Lender{ID: 3, IsVerified: true, TotalEarnings: 300.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM lenders WHERE id=$1)")).
					WithArgs(3).WillReturnError(errors.New("db failure"))
			},
			wantErr: true,
		},
		{
			name:   "insert fails",
			lender: &models.Lender{ID: 4, IsVerified: false, TotalEarnings: 400.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM lenders WHERE id=$1)")).
					WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
				mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO lenders (id, is_verified, total_earnings)
    VALUES ($1, $2, $3)
    RETURNING id
    `)).
					WithArgs(4, false, 400).WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()
			tt.mockSetup(mock)

			repo := lender_repo.NewLenderDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
			err := repo.Create(tt.lender)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLenderDBRepo_Save(t *testing.T) {
	sqlDB, _, _ := sqlmock.New()
	defer sqlDB.Close()

	repo := lender_repo.NewLenderDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
	if err := repo.Save(); err != nil {
		t.Errorf("Save() expected no error, got %v", err)
	}
}
