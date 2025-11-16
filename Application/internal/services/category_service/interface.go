package category_service

import "loopit/internal/models"

type CategoryServiceInterface interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(name string, price, security float64) error
	UpdateCategory(id int64, name string, price, security float64) error
	DeleteCategory(id int64) error
}
