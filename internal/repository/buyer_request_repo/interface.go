package buyer_request_repo

import "loopit/internal/models"

type BuyerRequestRepo interface {
	GetAllBuyerRequests(filterStatus []string) ([]models.BuyingRequest, error)
	UpdateStatusBuyerRequest(id int, newStatus string) error
	CreateBuyerRequest(req models.BuyingRequest) error
	GetBuyerRequestByID(id int) (*models.BuyingRequest, error)
	Save() error
}
