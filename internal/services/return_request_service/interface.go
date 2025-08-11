package return_request_service

import (
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
)

type ReturnRequestServiceInterface interface {
	CreateReturnRequest(userID int, orderID int) error
	UpdateReturnRequestStatus(userID int, reqID int, newStatus return_request_status.Status) error
	GetPendingReturnRequests(userID int) ([]models.ReturnRequest, error)
}
