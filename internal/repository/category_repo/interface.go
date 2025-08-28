//go:generate mockgen -source=interface.go -destination=../../mock/mock_category_repo.go -package=mock
package category_repo

import "loopit/internal/models"

type CategoryRepo interface {
	FindAll() ([]models.Category, error)
	FindByID(id int) (models.Category, error)
	Create(category models.Category) (err error)
	Save() error
}
