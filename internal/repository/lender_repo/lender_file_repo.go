// internal/repository/lender_file_repo.go
package lender_repo

import (
	"errors"
	"loopit/internal/models"
	"loopit/internal/storage"
)

type LenderFileRepo struct {
	lenderFile string
	lenders    []models.Lender
}

func NewLenderFileRepo(lenderFile string) (*LenderFileRepo, error) {
	lenders, err := storage.ReadJSONFile[models.Lender](lenderFile)
	if err != nil {
		return nil, err
	}

	return &LenderFileRepo{
		lenderFile: lenderFile,
		lenders:    lenders,
	}, nil
}

func (r *LenderFileRepo) FindAll() ([]models.Lender, error) {
	return r.lenders, nil
}

func (r *LenderFileRepo) FindByID(userID int) (*models.Lender, error) {
	for _, lender := range r.lenders {
		if lender.ID == userID {
			return &lender, nil
		}
	}
	return nil, errors.New("lender not found")
}

func (r *LenderFileRepo) Create(lender *models.Lender) error {
	// Avoid duplicate entry
	for _, l := range r.lenders {
		if l.ID == lender.ID {
			return nil // already exists, silently ignore
		}
	}
	r.lenders = append(r.lenders, *lender)
	return nil
}

func (r *LenderFileRepo) Save() error {
	return storage.WriteJSONFile(r.lenderFile, r.lenders)
}
