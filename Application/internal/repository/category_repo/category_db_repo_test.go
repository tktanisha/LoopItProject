package category_repo_test

// import (
// 	"database/sql"
// 	"errors"
// 	"regexp"
// 	"testing"

// 	"loopit/internal/db"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/category_repo"
// 	"loopit/pkg/logger"

// 	"github.com/DATA-DOG/go-sqlmock"
// )

// func TestCategoryDBRepo_FindAll(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 		wantCount int
// 	}{
// 		{
// 			name: "success with multiple categories",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "name", "price", "security"}).
// 					AddRow(1, "Electronics", 100.0, 50.0).
// 					AddRow(2, "Furniture", 200.0, 75.0)
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, security FROM categories")).
// 					WillReturnRows(rows)
// 			},
// 			wantErr:   false,
// 			wantCount: 2,
// 		},
// 		{
// 			name: "query error",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, security FROM categories")).
// 					WillReturnError(errors.New("db failure"))
// 			},
// 			wantErr:   true,
// 			wantCount: 0,
// 		},
// 		{
// 			name: "scan error on one row (skips row)",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"id", "name", "price", "security"}).
// 					AddRow(nil, "Invalid", 0.0, 0.0) // will cause scan error
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, security FROM categories")).
// 					WillReturnRows(rows)
// 			},
// 			wantErr:   false,
// 			wantCount: 0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			repo := category_repo.NewCategoryDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 			categories, err := repo.FindAll()

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
// 			}
// 			if len(categories) != tt.wantCount {
// 				t.Errorf("expected %d categories, got %d", tt.wantCount, len(categories))
// 			}
// 		})
// 	}
// }

// func TestCategoryDBRepo_FindByID(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 		expectNil bool
// 	}{
// 		{
// 			name: "success - category found",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				row := sqlmock.NewRows([]string{"id", "name", "price", "security"}).
// 					AddRow(1, "Electronics", 100.0, 50.0)
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, security FROM categories WHERE id=$1")).
// 					WithArgs(1).
// 					WillReturnRows(row)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "category not found",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, security FROM categories WHERE id=$1")).
// 					WithArgs(999).
// 					WillReturnError(sql.ErrNoRows)
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "db scan error",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, security FROM categories WHERE id=$1")).
// 					WithArgs(2).
// 					WillReturnError(errors.New("db scan failed"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			repo := category_repo.NewCategoryDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 			_, err := repo.FindByID(1)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }

// func TestCategoryDBRepo_Create(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		category  models.Category
// 		mockSetup func(sqlmock.Sqlmock)
// 		wantErr   bool
// 	}{
// 		{
// 			name:     "success - category created",
// 			category: models.Category{Name: "Electronics", Price: 100.0, Security: 50.0},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
//     INSERT INTO categories (name, price, security)
//     VALUES ($1, $2, $3)
//     RETURNING id
//     `)).
// 					WithArgs("Electronics", 100.0, 50.0).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:     "db error - category not created",
// 			category: models.Category{Name: "FailedCat", Price: 0.0, Security: 0.0},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(regexp.QuoteMeta(`
//     INSERT INTO categories (name, price, security)
//     VALUES ($1, $2, $3)
//     RETURNING id
//     `)).
// 					WithArgs("FailedCat", 0.0, 0.0).
// 					WillReturnError(errors.New("insert failed"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlDB, mock, _ := sqlmock.New()
// 			defer sqlDB.Close()
// 			tt.mockSetup(mock)

// 			repo := category_repo.NewCategoryDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 			err := repo.Create(tt.category)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }

// func TestCategoryDBRepo_Save(t *testing.T) {
// 	sqlDB, _, _ := sqlmock.New()
// 	defer sqlDB.Close()

// 	repo := category_repo.NewCategoryDBRepo(db.NewRealDatabase(sqlDB), logger.NewFakeLogger())
// 	if err := repo.Save(); err != nil {
// 		t.Errorf("expected no error from Save(), got %v", err)
// 	}
// }
