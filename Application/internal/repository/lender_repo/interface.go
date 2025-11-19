//go:generate mockgen -source=interface.go -destination=../../mock/mock_lender_repo.go -package=mock
package lender_repo

import "loopit/internal/models"

type LenderRepo interface {
	
	 Create(lender *models.Lender) error
	 Save() error
}
