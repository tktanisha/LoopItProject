package product_repo_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"loopit/internal/db"
	custom_mock "loopit/internal/mock"
	"loopit/internal/models"
	"loopit/internal/repository/product_repo"
	"loopit/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestProductDBRepo_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock, *custom_mock.MockCategoryRepo, *custom_mock.MockUserRepo)
		wantErr   bool
	}{
		{
			name: "success - products returned",
			mockSetup: func(mock sqlmock.Sqlmock, mockCat *custom_mock.MockCategoryRepo, mockUser *custom_mock.MockUserRepo) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products")).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "lender_id", "category_id", "name", "description", "duration", "is_available", "created_at",
					}).AddRow(1, 10, 100, "Product A", "Desc", 30, true, time.Now()))
				mockCat.EXPECT().FindByID(100).Return(models.Category{ID: 100, Name: "Category A"}, nil).Times(1)
				mockUser.EXPECT().FindByID(10).Return(&models.User{ID: 10, FullName: "Lender A"}, nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock, mockCat *custom_mock.MockCategoryRepo, mockUser *custom_mock.MockUserRepo) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products")).
					WillReturnError(errors.New("db failure"))
			},
			wantErr: true,
		},
		{
			name: "scan error - skip row",
			mockSetup: func(mock sqlmock.Sqlmock, mockCat *custom_mock.MockCategoryRepo, mockUser *custom_mock.MockUserRepo) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products")).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "lender_id", "category_id", "name", "description", "duration", "is_available", "created_at",
					}).AddRow("invalid", 10, 100, "Product A", "Desc", 30, true, time.Now()))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()

			mockCat := custom_mock.NewMockCategoryRepo(ctrl)
			mockUser := custom_mock.NewMockUserRepo(ctrl)

			tt.mockSetup(mock, mockCat, mockUser)

			repo := product_repo.NewProductDBRepo(db.NewRealDatabase(sqlDB), mockCat, mockUser, logger.NewFakeLogger())
			_, err := repo.FindAll()

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductDBRepo_FindByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(sqlmock.Sqlmock, *custom_mock.MockCategoryRepo, *custom_mock.MockUserRepo)
		wantErr   bool
	}{
		{
			name: "success - product found",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock, mockCat *custom_mock.MockCategoryRepo, mockUser *custom_mock.MockUserRepo) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products WHERE id=$1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "lender_id", "category_id", "name", "description", "duration", "is_available", "created_at",
					}).AddRow(1, 10, 100, "Product A", "Desc", 30, true, time.Now()))
				mockCat.EXPECT().FindByID(100).Return(models.Category{ID: 100, Name: "Category A"}, nil).Times(1)
				mockUser.EXPECT().FindByID(10).Return(&models.User{ID: 10, FullName: "Lender A"}, nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "no rows",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock, mockCat *custom_mock.MockCategoryRepo, mockUser *custom_mock.MockUserRepo) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products WHERE id=$1")).
					WithArgs(2).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "scan error",
			id:   3,
			mockSetup: func(mock sqlmock.Sqlmock, mockCat *custom_mock.MockCategoryRepo, mockUser *custom_mock.MockUserRepo) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products WHERE id=$1")).
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "lender_id", "category_id", "name", "description", "duration", "is_available", "created_at",
					}).AddRow("invalid", 10, 100, "Product A", "Desc", 30, true, time.Now()))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()

			mockCat := custom_mock.NewMockCategoryRepo(ctrl)
			mockUser := custom_mock.NewMockUserRepo(ctrl)

			tt.mockSetup(mock, mockCat, mockUser)

			repo := product_repo.NewProductDBRepo(db.NewRealDatabase(sqlDB), mockCat, mockUser, logger.NewFakeLogger())
			_, err := repo.FindByID(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductDBRepo_Create(t *testing.T) {
	tests := []struct {
		name      string
		product   *models.Product
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success - product created",
			product: &models.Product{
				LenderID:    10,
				CategoryID:  100,
				Name:        "Product A",
				Description: "Desc",
				Duration:    30,
				IsAvailable: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO products (lender_id, category_id, name, description, duration, is_available, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id
    `)).
					WithArgs(10, 100, "Product A", "Desc", 30, true, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "insert fails",
			product: &models.Product{
				LenderID:    11,
				CategoryID:  101,
				Name:        "Product B",
				Description: "Desc B",
				Duration:    15,
				IsAvailable: false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO products (lender_id, category_id, name, description, duration, is_available, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id
    `)).
					WithArgs(11, 101, "Product B", "Desc B", 15, false, sqlmock.AnyArg()).
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, _ := sqlmock.New()
			defer sqlDB.Close()

			tt.mockSetup(mock)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCat := custom_mock.NewMockCategoryRepo(ctrl)
			mockUser := custom_mock.NewMockUserRepo(ctrl)

			repo := product_repo.NewProductDBRepo(db.NewRealDatabase(sqlDB), mockCat, mockUser, logger.NewFakeLogger())
			err := repo.Create(tt.product)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductDBRepo_Save(t *testing.T) {
	repo := product_repo.NewProductDBRepo(nil, nil, nil, nil)
	if err := repo.Save(); err != nil {
		t.Errorf("Save() expected no error, got %v", err)
	}
}
