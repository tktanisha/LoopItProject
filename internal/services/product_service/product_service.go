package product_service

import (
	"errors"
	"fmt"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"
	"time"
)

type ProductService struct {
	productRepo product_repo.ProductRepo
	userRepo    user_repo.UserRepo
	log         *logger.Logger
}

func NewProductService(repo product_repo.ProductRepo, userRepo user_repo.UserRepo, log *logger.Logger) ProductServiceInterface {
	return &ProductService{productRepo: repo, userRepo: userRepo, log: log}
}

// GetAllProducts returns all products
func (p *ProductService) GetAllProducts() ([]*models.ProductResponse, error) {
	p.log.Info("Fetching all products")

	products, err := p.productRepo.FindAll()
	if err != nil {
		p.log.Error(fmt.Sprintf("Failed to fetch products: %v", err))
		return nil, err
	}

	p.log.Info(fmt.Sprintf("Fetched %d products successfully", len(products)))
	return products, nil
}

// GetProductByID returns a product by ID
func (p *ProductService) GetProductByID(id int) (*models.ProductResponse, error) {
	p.log.Info(fmt.Sprintf("Fetching product by ID: %d", id))

	if id <= 0 {
		p.log.Warning(fmt.Sprintf("Invalid product ID: %d", id))
		return nil, errors.New("product ID must be a positive integer")
	}

	product, err := p.productRepo.FindByID(id)
	if err != nil {
		p.log.Error(fmt.Sprintf("Product not found for ID %d: %v", id, err))
		return nil, errors.New("product not found")
	}

	p.log.Info(fmt.Sprintf("Fetched product successfully: ID %d", id))
	return product, nil
}

// CreateProduct creates a new product with default values
func (p *ProductService) CreateProduct(product *models.Product, userCtx *models.UserContext) error {
	if userCtx == nil {
		p.log.Error("Attempt to create product without user context")
		return errors.New("user not logged in")
	}

	p.log.Info(fmt.Sprintf("Creating product by user ID %d", userCtx.ID))

	if userCtx.Role != enums.RoleLender {
		p.log.Warning(fmt.Sprintf("Unauthorized product creation attempt by user ID %d with role %s", userCtx.ID, userCtx.Role))
		return errors.New("only lenders can create products")
	}

	product.LenderID = userCtx.ID
	product.CreatedAt = time.Now()

	err := p.productRepo.Create(product)
	if err != nil {
		p.log.Error(fmt.Sprintf("Failed to create product for user ID %d: %v", userCtx.ID, err))
		return err
	}

	p.log.Info(fmt.Sprintf("Product created successfully with ID %d by user %d", product.ID, userCtx.ID))
	return nil
}
