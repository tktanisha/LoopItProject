package product_repo

import (
	"errors"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/storage"
)

type ProductFileRepo struct {
	productFile  string
	products     []models.Product
	categoryRepo category_repo.CategoryRepo
	userRepo     user_repo.UserRepo
}

func NewProductFileRepo(productFile string, categoryRepo category_repo.CategoryRepo, userRepo user_repo.UserRepo) (*ProductFileRepo, error) {
	products, err := storage.ReadJSONFile[models.Product](productFile)
	if err != nil {
		return nil, err
	}
	return &ProductFileRepo{
		productFile:  productFile,
		products:     products,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}, nil
}

func (r *ProductFileRepo) FindAll() ([]*models.ProductResponse, error) {
	var responses []*models.ProductResponse

	for _, product := range r.products {
		category, _ := r.categoryRepo.FindByID(product.CategoryID)
		user, _ := r.userRepo.FindByID(product.LenderID)

		responses = append(responses, &models.ProductResponse{
			Product:  product,
			Category: category,
			User:     *user,
		})
	}

	return responses, nil
}

func (r *ProductFileRepo) FindByID(id int) (*models.ProductResponse, error) {
	for _, product := range r.products {
		if product.ID == id {
			category, _ := r.categoryRepo.FindByID(product.CategoryID)
			user, _ := r.userRepo.FindByID(product.LenderID)

			return &models.ProductResponse{
				Product:  product,
				Category: category,
				User:     *user,
			}, nil
		}
	}
	return nil, errors.New("product not found")
}

func (r *ProductFileRepo) Create(product *models.Product) error {
	product.ID = len(r.products) + 1
	r.products = append(r.products, *product)
	return nil
}

func (r *ProductFileRepo) Save() error {
	return storage.WriteJSONFile(r.productFile, r.products)
}
