package category_service

import (
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
)

type CategoryService struct {
	categoryRepo category_repo.CategoryRepo
}

func NewCategoryService(repo category_repo.CategoryRepo) CategoryServiceInterface {
	return &CategoryService{
		categoryRepo: repo,
	}
}

func (c *CategoryService) GetAllCategories() ([]models.Category, error) {
	categories, err := c.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *CategoryService) CreateCategory(name string, price, security float64) error {
	category := models.Category{
		Name:     name,
		Price:    price,
		Security: security,
	}
	if err := c.categoryRepo.Create(category); err != nil {
		return err
	}
	return nil
}

func (c *CategoryService) UpdateCategory(id int64, name string, price, security float64) error {
	category, err := c.categoryRepo.FindByID(id)
	
	if err != nil {
		return err
	}
	category.Name = name
	category.Price = price
	category.Security = security

	if err := c.categoryRepo.Update(category); err != nil {
		return err
	}
	return nil
}

func (c *CategoryService) DeleteCategory(id int64) error {
	if err := c.categoryRepo.Delete(id); err != nil {
		return err
	}
	return nil
}
