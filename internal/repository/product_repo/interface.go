//go:generate mockgen -source=interface.go -destination=../../mock/mock_product_repo.go -package=mock
package product_repo

import "loopit/internal/models"

type ProductRepo interface {
	FindAll(models.ProductFilter) ([]*models.ProductResponse, error)
	FindByID(id int) (*models.ProductResponse, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id int) error
	Save() error
}
