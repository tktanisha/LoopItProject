// file: internal/services/category_service/category_service_test.go
package category_service_test

// import (
// 	"errors"
// 	"loopit/internal/models"
// 	"loopit/internal/services/category_service"
// 	"loopit/pkg/logger"
// 	"testing"
// )

// // FakeRepo implements category_repo.CategoryRepo
// type FakeRepo struct {
// 	FindAllFn func() ([]models.Category, error)
// 	CreateFn  func(category models.Category) error
// }

// func (f *FakeRepo) FindAll() ([]models.Category, error) {
// 	return f.FindAllFn()
// }
// func (f *FakeRepo) FindByID(id int) (models.Category, error) {
// 	return models.Category{}, nil
// }
// func (f *FakeRepo) Create(category models.Category) error {
// 	return f.CreateFn(category)
// }
// func (f *FakeRepo) Save() error { return nil }

// func TestGetAllCategories(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name    string
// 		repoFn  func() ([]models.Category, error)
// 		wantErr bool
// 	}{
// 		{
// 			name: "success",
// 			repoFn: func() ([]models.Category, error) {
// 				return []models.Category{
// 					{ID: 1, Name: "Electronics"},
// 				}, nil
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repo error",
// 			repoFn: func() ([]models.Category, error) {
// 				return nil, errors.New("db error")
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			repo := &FakeRepo{FindAllFn: tt.repoFn}
// 			svc := category_service.NewCategoryService(repo, log)
// 			_, err := svc.GetAllCategories()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("expected err=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }

// func TestCreateCategory(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name    string
// 		repoFn  func(models.Category) error
// 		wantErr bool
// 	}{
// 		{
// 			name: "success",
// 			repoFn: func(c models.Category) error {
// 				return nil
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repo error",
// 			repoFn: func(c models.Category) error {
// 				return errors.New("insert failed")
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			repo := &FakeRepo{CreateFn: tt.repoFn}
// 			svc := category_service.NewCategoryService(repo, log)
// 			err := svc.CreateCategory("Books", 100, 50)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("expected err=%v, got %v", tt.wantErr, err)
// 			}
// 		})
// 	}
// }
