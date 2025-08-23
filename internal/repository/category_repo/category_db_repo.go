package category_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/models"
	"loopit/pkg/logger"
)

type CategoryDBRepo struct {
	db  *sql.DB
	log *logger.Logger
}

func NewCategoryDBRepo(db *sql.DB, log *logger.Logger) *CategoryDBRepo {
	return &CategoryDBRepo{db: db, log: log}
}

// FindAll fetches all categories from the database
func (r *CategoryDBRepo) FindAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, price, security FROM categories")
	if err != nil {
		if r.log != nil {
			r.log.Error(fmt.Sprintf("Repo: DB error fetching all categories: %v", err))
		}
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Price, &c.Security); err != nil {
			if r.log != nil {
				r.log.Warning(fmt.Sprintf("Repo: Could not scan category row: %v", err))
			}
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// FindByID fetches a category by its ID
func (r *CategoryDBRepo) FindByID(id int) (models.Category, error) {
	row := r.db.QueryRow("SELECT id, name, price, security FROM categories WHERE id=$1", id)
	var c models.Category
	if err := row.Scan(&c.ID, &c.Name, &c.Price, &c.Security); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if r.log != nil {
				r.log.Warning(fmt.Sprintf("Repo: No category found in DB with id %d", id))
			}
			return models.Category{}, errors.New("category not found")
		}
		if r.log != nil {
			r.log.Error(fmt.Sprintf("Repo: DB error scanning category by id %d: %v", id, err))
		}
		return models.Category{}, err
	}
	return c, nil
}

// Create inserts a new category into the database
func (r *CategoryDBRepo) Create(category models.Category) error {
	query := `
    INSERT INTO categories (name, price, security)
    VALUES ($1, $2, $3)
    RETURNING id
    `
	err := r.db.QueryRow(query, category.Name, category.Price, category.Security).Scan(&category.ID)
	if err != nil && r.log != nil {
		r.log.Error(fmt.Sprintf("Repo: DB error creating category '%s': %v", category.Name, err))
	}
	return err
}

// Save is a no-op for Postgres
func (r *CategoryDBRepo) Save() error {
	return nil
}
