package lender_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/models"
	"loopit/pkg/logger"
)

type LenderDBRepo struct {
	db  *sql.DB
	log *logger.Logger
}

func NewLenderDBRepo(db *sql.DB, log *logger.Logger) *LenderDBRepo {
	return &LenderDBRepo{db: db, log: log}
}

// FindAll returns all lenders
func (r *LenderDBRepo) FindAll() ([]models.Lender, error) {
	rows, err := r.db.Query("SELECT id, is_verified, total_earnings FROM lenders")
	if err != nil {
		if r.log != nil {
			r.log.Error(fmt.Sprintf("DB error fetching lenders: %v", err))
		}
		return nil, err
	}
	defer rows.Close()

	var lenders []models.Lender
	for rows.Next() {
		var l models.Lender
		if err := rows.Scan(&l.ID, &l.IsVerified, &l.TotalEarnings); err != nil {
			if r.log != nil {
				r.log.Warning(fmt.Sprintf("DB warning scanning lender row: %v", err))
			}
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
			if r.log != nil {
				r.log.Warning(fmt.Sprintf("DB: No lender found with id %d", userID))
			}
			return nil, errors.New("lender not found")
		}
		if r.log != nil {
			r.log.Error(fmt.Sprintf("DB error fetching lender by id %d: %v", userID, err))
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
		if r.log != nil {
			r.log.Error(fmt.Sprintf("DB error checking lender existence for id %d: %v", lender.ID, err))
		}
		return err
	}
	if exists {
		if r.log != nil {
			r.log.Info(fmt.Sprintf("Lender with id %d already exists, skipping create", lender.ID))
		}
		return nil // already exists
	}

	query := `
    INSERT INTO lenders (id, is_verified, total_earnings)
    VALUES ($1, $2, $3)
    RETURNING id
    `
	err = r.db.QueryRow(query, lender.ID, lender.IsVerified, lender.TotalEarnings).Scan(&lender.ID)
	if err != nil && r.log != nil {
		r.log.Error(fmt.Sprintf("DB error inserting lender with id %d: %v", lender.ID, err))
	}
	return err
}

// Save is a no-op for Postgres
func (r *LenderDBRepo) Save() error {
	return nil
}
