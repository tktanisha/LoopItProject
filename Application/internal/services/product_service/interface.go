package product_service

import "loopit/internal/models"

type ProductServiceInterface interface {
	GetAllProducts(models.ProductFilter) ([]*models.ProductResponse, error)
	GetProductByID(id int64) (*models.ProductResponse, error)
	CreateProduct(product *models.Product, userCtx *models.UserContext) error
	UpdateProduct(productID int64, name string, description string, categoryID int64, userCtx *models.UserContext) error
	DeleteProduct(id int64, userCtx *models.UserContext) error
}