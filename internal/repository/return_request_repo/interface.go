package return_request_repo

import "loopit/internal/models"

type ReturnRequestRepo interface {
	CreateReturnRequest(req models.ReturnRequest) error
	UpdateReturnRequestStatus(id int, newStatus string) error
	GetAllReturnRequests(filterStatuses []string) ([]models.ReturnRequest, error)
	GetReturnRequestByID(id int) (models.ReturnRequest, error)
	Save() error
}
