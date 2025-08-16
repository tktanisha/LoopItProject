package lender_repo

import (
	"database/sql"
	"errors"
	"loopit/internal/models"
)

type LenderDBRepo struct {
	db *sql.DB
}

func NewLenderDBRepo(db *sql.DB) *LenderDBRepo {
	return &LenderDBRepo{db: db}
}

// FindAll returns all lenders
func (r *LenderDBRepo) FindAll() ([]models.Lender, error) {
	rows, err := r.db.Query("SELECT id, is_verified, total_earnings FROM lenders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lenders []models.Lender
	for rows.Next() {
		var l models.Lender
		if err := rows.Scan(&l.ID, &l.IsVerified, &l.TotalEarnings); err != nil {
			continue
		}
		lenders = append(lenders, l)
	}

	return lenders, nil
}

// FindByID returns a lender by user ID
func (r *LenderDBRepo) FindByID(userID int) (*models.Lender, error) {
	row := r.db.QueryRow("SELECT id, is_verified, total_earnings FROM lenders WHERE id=$1", userID)
	var l models.Lender
	if err := row.Scan(&l.ID, &l.IsVerified, &l.TotalEarnings); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("lender not found")
		}
		return nil, err
	}
	return &l, nil
}

// Create inserts a new lender into the database
func (r *LenderDBRepo) Create(lender *models.Lender) error {
	// Check if lender already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM lenders WHERE id=$1)", lender.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // already exists
	}

	query := `
	INSERT INTO lenders (id, is_verified, total_earnings)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	return r.db.QueryRow(query, lender.ID, lender.IsVerified, lender.TotalEarnings).Scan(&lender.ID)
}

// Save is a no-op for Postgres
func (r *LenderDBRepo) Save() error {
	return nil
}

// // Optional: CreateWithTx allows transactional creation (used for BecomeLender in UserDBRepo)
// func (r *LenderDBRepo) CreateWithTx(tx *sql.Tx, lender *models.Lender) error {
// 	var exists bool
// 	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM lenders WHERE id=$1)", lender.ID).Scan(&exists)
// 	if err != nil {
// 		return err
// 	}
// 	if exists {
// 		return nil
// 	}

// 	query := `
// 	INSERT INTO lenders (id, is_verified, total_earnings)
// 	VALUES ($1, $2, $3)
// 	RETURNING id
// 	`
// 	return tx.QueryRow(query, lender.ID, lender.IsVerified, lender.TotalEarnings).Scan(&lender.ID)
// }
