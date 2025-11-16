package product_service

import (
	"errors"
	"fmt"
	"loopit/internal/constants"
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
func (p *ProductService) GetAllProducts(filters models.ProductFilter) ([]*models.ProductResponse, error) {

	products, err := p.productRepo.FindAll(filters)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductByID returns a product by ID
func (p *ProductService) GetProductByID(id int64) (*models.ProductResponse, error) {
	fmt.Println("the product of id=", id)

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
	product.IsAvailable = true

	product.ImageUrl = constants.CategoryImageMap[fmt.Sprint(product.CategoryID)]
	if product.ImageUrl == "" {
		product.ImageUrl = constants.CategoryImageMap["Default"]
	}

	err := p.productRepo.Create(product)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct updates an existing product
func (p *ProductService) UpdateProduct(productID int64, name string, description string, categoryID int64, userCtx *models.UserContext) error {
	if userCtx == nil {
		return errors.New("user not logged in")
	}
	if userCtx.Role != enums.RoleLender {
		return errors.New("only lenders can update products")
	}
	product, err := p.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("product not found")
	}
	if product.Product.LenderID != userCtx.ID {
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
		return err
	}
	return nil
}

// DeleteProduct deletes a product by ID
func (p *ProductService) DeleteProduct(id int64, userCtx *models.UserContext) error {
	if userCtx == nil {
		return errors.New("user not logged in")
	}

	if userCtx.Role != enums.RoleLender {
		return errors.New("only lenders can delete products")
	}
	product, err := p.productRepo.FindByID(id)

	if err != nil {
		return errors.New("product not found")

	}
	if product.Product.LenderID != userCtx.ID {

		return errors.New("you can only delete your own products")
	}
	err = p.productRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
