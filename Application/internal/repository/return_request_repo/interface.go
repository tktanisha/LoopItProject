//go:generate mockgen -source=interface.go -destination=../../mock/mock_return_request_repo.go -package=mock
package return_request_repo

import "loopit/internal/models"

type ReturnRequestRepo interface {
	CreateReturnRequest(req models.ReturnRequest) error
	UpdateReturnRequestStatus(id int64, newStatus string) error
	GetAllReturnRequests(filterStatuses []string) ([]models.ReturnRequest, error)
	GetReturnRequestByID(id int64) (models.ReturnRequest, error)
	Save() error
}
