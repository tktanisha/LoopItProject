package product_service

import "loopit/internal/models"

type ProductServiceInterface interface {
	GetAllProducts(models.ProductFilter) ([]*models.ProductResponse, error)
	GetProductByID(id int) (*models.ProductResponse, error)
	CreateProduct(product *models.Product, userCtx *models.UserContext) error
	UpdateProduct(productID int, name string, description string, price float64, categoryID int, userCtx *models.UserContext) error
	DeleteProduct(id int, userCtx *models.UserContext) error
}
