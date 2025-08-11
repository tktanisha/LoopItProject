package category_repo

import "loopit/internal/models"

type CategoryRepo interface {
	FindAll() ([]models.Category, error)
	FindByID(id int) (models.Category, error)
	Create(category models.Category)
	Save() error
}
