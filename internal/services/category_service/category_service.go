package category_service

import (
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
)

type CategoryService struct {
	categoryRepo category_repo.CategoryRepo
}

func NewCategoryService(repo category_repo.CategoryRepo) CategoryServiceInterface {
	return &CategoryService{categoryRepo: repo}
}

func (c *CategoryService) GetAllCategories() ([]models.Category, error) {
	return c.categoryRepo.FindAll()
}

func (c *CategoryService) CreateCategory(name string, price, security float64) error {
	category := models.Category{
		Name:     name,
		Price:    price,
		Security: security,
	}
	c.categoryRepo.Create(category)
	return nil
}
