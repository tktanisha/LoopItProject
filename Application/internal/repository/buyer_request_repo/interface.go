//go:generate mockgen -source=interface.go -destination=../../mock/mock_buyer_request_repo.go -package=mock
package buyer_request_repo

import "loopit/internal/models"

type BuyerRequestRepo interface {
	GetAllBuyerRequests(id *int64, filterStatus []string) ([]models.BuyingRequest, error)
	UpdateStatusBuyerRequest(id int64, newStatus string) error
	CreateBuyerRequest(req models.BuyingRequest) error
	GetBuyerRequestByID(id int64) (*models.BuyingRequest, error)
	Save() error
}
