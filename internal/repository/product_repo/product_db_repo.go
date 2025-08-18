package product_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"
	"time"
)

type ProductDBRepo struct {
	db           *sql.DB
	categoryRepo category_repo.CategoryRepo
	userRepo     user_repo.UserRepo
	log          *logger.Logger
}

func NewProductDBRepo(db *sql.DB, categoryRepo category_repo.CategoryRepo, userRepo user_repo.UserRepo, log *logger.Logger) *ProductDBRepo {
	return &ProductDBRepo{
		db:           db,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		log:          log,
	}
}

// FindAll returns all products with category and lender info
func (r *ProductDBRepo) FindAll() ([]*models.ProductResponse, error) {
	rows, err := r.db.Query("SELECT id, lender_id, category_id, name, description, duration, is_available, created_at FROM products")
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
		if err := rows.Scan(&p.ID, &p.LenderID, &p.CategoryID, &p.Name, &p.Description, &p.Duration, &p.IsAvailable, &p.CreatedAt); err != nil {
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
func (r *ProductDBRepo) FindByID(id int) (*models.ProductResponse, error) {
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
    INSERT INTO products (lender_id, category_id, name, description, duration, is_available, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id
    `
	err := r.db.QueryRow(query, product.LenderID, product.CategoryID, product.Name, product.Description, product.Duration, product.IsAvailable, time.Now()).Scan(&product.ID)
	if err != nil && r.log != nil {
		r.log.Error(fmt.Sprintf("DB error creating product '%s': %v", product.Name, err))
	}
	return err
}

// Save is a no-op for Postgres as changes are applied immediately
func (r *ProductDBRepo) Save() error {
	return nil
}
