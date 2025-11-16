package product_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"
	"time"
)

type ProductDBRepo struct {
	db           *db.DynamoClient
	categoryRepo category_repo.CategoryRepo
	userRepo     user_repo.UserRepo
	log          logger.LoggerInterface
}

func NewProductDBRepo(db *db.DynamoClient, categoryRepo category_repo.CategoryRepo, userRepo user_repo.UserRepo) *ProductDBRepo {
	return &ProductDBRepo{
		db:           db,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

// FindAll returns all products with category and lender info
func (r *ProductDBRepo) FindAll(filters models.ProductFilter) ([]*models.ProductResponse, error) {
	query := `
		SELECT id, lender_id, category_id, name, description, duration, is_available, created_at, image_url
		FROM products WHERE 1=1
	`
	args := []any{}
	argIdx := 1 // For PostgreSQL placeholders

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (LOWER(name) LIKE LOWER($%d) OR LOWER(description) LIKE LOWER($%d))", argIdx, argIdx+1)
		searchParam := "%" + filters.Search + "%"
		args = append(args, searchParam, searchParam)
		argIdx += 2
	}
	if filters.LenderID != "" {
		query += fmt.Sprintf(" AND lender_id = $%d", argIdx)
		args = append(args, filters.LenderID)
		argIdx++
	}
	if filters.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argIdx)
		args = append(args, filters.CategoryID)
		argIdx++
	}
	if filters.IsAvailable != "" {
		query += fmt.Sprintf(" AND is_available = $%d", argIdx)
		args = append(args, filters.IsAvailable)
		argIdx++
	}

	r.log.Info(fmt.Sprintf("Executing product query: %s with args %v", query, args))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		if r.log != nil {
			r.log.Error(fmt.Sprintf("DB error fetching products: %v", err))
		}
		return nil, err
	}
	defer rows.Close()

	var responses []*models.ProductResponse
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.LenderID, &p.CategoryID, &p.Name, &p.Description, &p.Duration, &p.IsAvailable, &p.CreatedAt, &p.ImageUrl); err != nil {
			if r.log != nil {
				r.log.Warning(fmt.Sprintf("DB warning scanning product row: %v", err))
			}
			continue
		}

		category, _ := r.categoryRepo.FindByID(p.CategoryID)
		user, _ := r.userRepo.FindByID(p.LenderID)

		responses = append(responses, &models.ProductResponse{
			Product:  p,
			Category: category,
			User:     *user,
		})
	}

	return responses, nil
}

// FindByID returns a single product by ID with category and lender info
func (r *ProductDBRepo) FindByID(id int64) (*models.ProductResponse, error) {
	row := r.db.QueryRow("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products WHERE id=$1", id)

	var p models.Product
	if err := row.Scan(&p.ID, &p.LenderID, &p.CategoryID, &p.Name, &p.Description, &p.Duration, &p.IsAvailable, &p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if r.log != nil {
				r.log.Warning(fmt.Sprintf("DB: No product found with id %d", id))
			}
			return nil, errors.New("product not found")
		}
		if r.log != nil {
			r.log.Error(fmt.Sprintf("DB error fetching product by id %d: %v", id, err))
		}
		return nil, err
	}

	category, _ := r.categoryRepo.FindByID(p.CategoryID)
	user, _ := r.userRepo.FindByID(p.LenderID)

	return &models.ProductResponse{
		Product:  p,
		Category: category,
		User:     *user,
	}, nil
}

// Create inserts a new product into the database
func (r *ProductDBRepo) Create(product *models.Product) error {
	query := `
    INSERT INTO products (lender_id, category_id, name, description, duration, is_available, image_url, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id
    `
	err := r.db.QueryRow(query, product.LenderID, product.CategoryID, product.Name, product.Description, product.Duration, product.IsAvailable, product.ImageUrl, time.Now()).Scan(&product.ID)
	if err != nil && r.log != nil {
		r.log.Error(fmt.Sprintf("DB error creating product '%s': %v", product.Name, err))
	}
	return err
}

// update
func (r *ProductDBRepo) Update(product *models.Product) error {
	query := `
	UPDATE products
	SET lender_id=$1, category_id=$2, name=$3, description=$4, duration=$5, is_available=$6
	WHERE id=$7
	`
	_, err := r.db.Exec(query, product.LenderID, product.CategoryID, product.Name, product.Description, product.Duration, product.IsAvailable, product.ID)
	if err != nil && r.log != nil {
		r.log.Error(fmt.Sprintf("DB error updating product ID %d: %v", product.ID, err))
	}
	return err
}

// delete
func (r *ProductDBRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil && r.log != nil {
		r.log.Error(fmt.Sprintf("DB error deleting product ID %d: %v", id, err))
	}
	return err
}

// Save is a no-op for Postgres as changes are applied immediately
func (r *ProductDBRepo) Save() error {
	return nil
}
