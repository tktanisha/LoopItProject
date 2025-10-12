package category_service

import (
	"fmt"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/pkg/logger"
)

type CategoryService struct {
	categoryRepo category_repo.CategoryRepo
	log          logger.LoggerInterface
}

func NewCategoryService(repo category_repo.CategoryRepo, log logger.LoggerInterface) CategoryServiceInterface {
	return &CategoryService{
		categoryRepo: repo,
		log:          log,
	}
}

func (c *CategoryService) GetAllCategories() ([]models.Category, error) {
	c.log.Info("Service: Fetching all categories")
	categories, err := c.categoryRepo.FindAll()
	if err != nil {
		c.log.Error(fmt.Sprintf("Service: Failed to fetch categories, error: %v", err))
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Service: Successfully fetched %d categories", len(categories)))
	return categories, nil
}

func (c *CategoryService) CreateCategory(name string, price, security float64) error {
	c.log.Info(fmt.Sprintf("Service: Creating category '%s' with price %.2f and security %.2f", name, price, security))
	category := models.Category{
		Name:     name,
		Price:    price,
		Security: security,
	}
	if err := c.categoryRepo.Create(category); err != nil {
		c.log.Error(fmt.Sprintf("Service: Failed to create category '%s', error: %v", name, err))
		return err
	}
	c.log.Info(fmt.Sprintf("Service: Category created successfully: '%s'", name))
	return nil
}

func (c *CategoryService) UpdateCategory(id int, name string, price, security float64) error {
	c.log.Info(fmt.Sprintf("Service: Updating category ID %d to name '%s', price %.2f, security %.2f", id, name, price, security))
	category, err := c.categoryRepo.FindByID(id)
	
	if err != nil {
		c.log.Error(fmt.Sprintf("Service: Failed to find category ID %d, error: %v", id, err))
		return err
	}
	category.Name = name
	category.Price = price
	category.Security = security

	if err := c.categoryRepo.Update(category); err != nil {
		c.log.Error(fmt.Sprintf("Service: Failed to update category ID %d, error: %v", id, err))
		return err
	}
	c.log.Info(fmt.Sprintf("Service: Category ID %d updated successfully", id))
	return nil
}

func (c *CategoryService) DeleteCategory(id int) error {
	c.log.Info(fmt.Sprintf("Service: Deleting category ID %d", id))
	if err := c.categoryRepo.Delete(id); err != nil {
		c.log.Error(fmt.Sprintf("Service: Failed to delete category ID %d, error: %v", id, err))
		return err
	}
	c.log.Info(fmt.Sprintf("Service: Category ID %d deleted successfully", id))
	return nil
}
