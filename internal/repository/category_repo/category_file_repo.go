package category_repo

import (
	"errors"
	"loopit/internal/models"
	"loopit/internal/storage"
)

type CategoryFileRepo struct {
	categoryFile string
	categories   []models.Category
}

func NewCategoryFileRepo(categoryFile string) (*CategoryFileRepo, error) {
	categories, err := storage.ReadJSONFile[models.Category](categoryFile)
	if err != nil {
		return nil, err
	}
	return &CategoryFileRepo{
		categoryFile: categoryFile,
		categories:   categories,
	}, nil
}

func (r *CategoryFileRepo) FindAll() ([]models.Category, error) {
	return r.categories, nil
}

func (r *CategoryFileRepo) FindByID(id int) (models.Category, error) {
	for _, c := range r.categories {
		if c.ID == id {
			return c, nil
		}
	}
	return models.Category{}, errors.New("category not found")
}

func (r *CategoryFileRepo) Create(category models.Category) {
	category.ID = len(r.categories) + 1
	r.categories = append(r.categories, category)
}

func (r *CategoryFileRepo) Save() error {
	return storage.WriteJSONFile(r.categoryFile, r.categories)
}
