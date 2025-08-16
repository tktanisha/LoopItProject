package category_repo

import (
	"database/sql"
	"errors"
	"loopit/internal/models"
)

type CategoryDBRepo struct {
	db *sql.DB
}

func NewCategoryDBRepo(db *sql.DB) *CategoryDBRepo {
	return &CategoryDBRepo{db: db}
}

// FindAll fetches all categories from the database
func (r *CategoryDBRepo) FindAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, price, security FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Price, &c.Security); err != nil {
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
			return models.Category{}, errors.New("category not found")
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
	return r.db.QueryRow(query, category.Name, category.Price, category.Security).Scan(category.ID)
}

// Save is a no-op for Postgres
func (r *CategoryDBRepo) Save() error {
	return nil
}
