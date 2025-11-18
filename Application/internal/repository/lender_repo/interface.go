//go:generate mockgen -source=interface.go -destination=../../mock/mock_lender_repo.go -package=mock
package lender_repo

type LenderRepo interface {
	// FindAll() ([]models.Lender, error)
	// FindByID(userID int64) (*models.Lender, error)
	// Create(lender *models.Lender) error
	 Save() error
}
