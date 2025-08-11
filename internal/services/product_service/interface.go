package product_service

import "loopit/internal/models"

type ProductServiceInterface interface {
	GetAllProducts() ([]*models.ProductResponse, error)
	GetProductByID(id int) (*models.ProductResponse, error)
	CreateProduct(product *models.Product, userCtx *models.UserContext) error
}
