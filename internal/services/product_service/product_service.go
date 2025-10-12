package product_service

import (
	"errors"
	"fmt"
	"loopit/internal/constants"
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
	log         logger.LoggerInterface
}

func NewProductService(repo product_repo.ProductRepo, userRepo user_repo.UserRepo, log logger.LoggerInterface) ProductServiceInterface {
	return &ProductService{productRepo: repo, userRepo: userRepo, log: log}
}

// GetAllProducts returns all products
func (p *ProductService) GetAllProducts(filters models.ProductFilter) ([]*models.ProductResponse, error) {
	p.log.Info("Fetching all products with filters")

	products, err := p.productRepo.FindAll(filters)
	if err != nil {
		p.log.Error(fmt.Sprintf("Failed to fetch products: %v", err))
		return nil, err
	}

	p.log.Info(fmt.Sprintf("Fetched %d products successfully", len(products)))
	return products, nil
}

// GetProductByID returns a product by ID
func (p *ProductService) GetProductByID(id int) (*models.ProductResponse, error) {
	fmt.Println("the product of id=", id)
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
	fmt.Println("products in service=", product)
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
	product.IsAvailable = true

	product.ImageUrl = constants.CategoryImageMap[fmt.Sprint(product.CategoryID)]
	if product.ImageUrl == "" {
		product.ImageUrl = constants.CategoryImageMap["default"]
	}

	err := p.productRepo.Create(product)
	if err != nil {
		p.log.Error(fmt.Sprintf("Failed to create product for user ID %d: %v", userCtx.ID, err))
		return err
	}

	p.log.Info(fmt.Sprintf("Product created successfully with ID %d by user %d", product.ID, userCtx.ID))
	return nil
}

// UpdateProduct updates an existing product
func (p *ProductService) UpdateProduct(productID int, name string, description string, categoryID int, userCtx *models.UserContext) error {
	if userCtx == nil {
		p.log.Error("Attempt to update product without user context")
		return errors.New("user not logged in")
	}
	p.log.Info(fmt.Sprintf("Updating product ID %d by user ID %d", productID, userCtx.ID))
	if userCtx.Role != enums.RoleLender {
		p.log.Warning(fmt.Sprintf("Unauthorized product update attempt by user ID %d with role %s", userCtx.ID, userCtx.Role))
		return errors.New("only lenders can update products")
	}
	product, err := p.productRepo.FindByID(productID)
	if err != nil {
		p.log.Error(fmt.Sprintf("Product not found for ID %d: %v", productID, err))
		return errors.New("product not found")
	}
	if product.Product.LenderID != userCtx.ID {
		p.log.Warning(fmt.Sprintf("User ID %d attempted to update product ID %d they do not own", userCtx.ID, productID))
		return errors.New("you can only update your own products")
	}
	product.Product.Name = name

	product.Product.Description = description
	product.Product.CategoryID = categoryID
	product.Product.ImageUrl = constants.CategoryImageMap[fmt.Sprint(categoryID)]
	if product.Product.ImageUrl == "" {
		product.Product.ImageUrl = constants.CategoryImageMap["default"]
	}
	err = p.productRepo.Update(&product.Product)
	if err != nil {
		p.log.Error(fmt.Sprintf("Failed to update product ID %d: %v", productID, err))
		return err
	}
	p.log.Info(fmt.Sprintf("Product ID %d updated successfully by user ID %d", productID, userCtx.ID))
	return nil
}

// DeleteProduct deletes a product by ID
func (p *ProductService) DeleteProduct(id int, userCtx *models.UserContext) error {
	if userCtx == nil {
		p.log.Error("Attempt to delete product without user context")
		return errors.New("user not logged in")
	}

	p.log.Info(fmt.Sprintf("Deleting product ID %d by user ID %d", id, userCtx.ID))
	if userCtx.Role != enums.RoleLender {
		p.log.Warning(fmt.Sprintf("Unauthorized product deletion attempt by user ID %d with role %s", userCtx.ID, userCtx.Role))
		return errors.New("only lenders can delete products")
	}
	product, err := p.productRepo.FindByID(id)

	if err != nil {
		p.log.Error(fmt.Sprintf("Product not found for ID %d: %v", id, err))
		return errors.New("product not found")

	}
	if product.Product.LenderID != userCtx.ID {
		p.log.Warning(fmt.Sprintf("User ID %d attempted to delete product ID %d they do not own", userCtx.ID, id))

		return errors.New("you can only delete your own products")
	}
	err = p.productRepo.Delete(id)
	if err != nil {
		p.log.Error(fmt.Sprintf("Failed to delete product ID %d: %v", id, err))
		return err
	}
	p.log.Info(fmt.Sprintf("Product ID %d deleted successfully by user ID %d", id, userCtx.ID))
	return nil
}
