package product_service

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/user_repo"
	"time"
)

type ProductService struct {
	productRepo product_repo.ProductRepo
	userRepo    user_repo.UserRepo
}

func NewProductService(repo product_repo.ProductRepo, userRepo user_repo.UserRepo) ProductServiceInterface {
	return &ProductService{productRepo: repo, userRepo: userRepo}
}

// GetAllProducts returns all products
func (p *ProductService) GetAllProducts() ([]*models.ProductResponse, error) {
	return p.productRepo.FindAll()
}

// GetProductByID returns a product by ID
func (p *ProductService) GetProductByID(id int) (*models.ProductResponse, error) {
	if id <= 0 {
		return nil, errors.New("product ID must be a positive integer")
	}

	product, err := p.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// CreateProduct creates a new product with default values
func (p *ProductService) CreateProduct(product *models.Product, userCtx *models.UserContext) error {
	if userCtx == nil {
		return errors.New("user not logged in")
	}

	if userCtx.Role != enums.RoleLender {
		return errors.New("only lenders can create products")
	}

	product.LenderID = userCtx.ID
	product.CreatedAt = time.Now()
	return p.productRepo.Create(product)
}
