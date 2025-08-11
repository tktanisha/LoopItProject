package product_repo

import "loopit/internal/models"

type ProductRepo interface {
	FindAll() ([]*models.ProductResponse, error)
	FindByID(id int) (*models.ProductResponse, error)
	Create(product *models.Product) error
	Save() error
}
