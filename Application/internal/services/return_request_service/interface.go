package return_request_service

import (
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
)

type ReturnRequestServiceInterface interface {
	CreateReturnRequest(userID int64, orderID int64) error
	UpdateReturnRequestStatus(userID int64, reqID int64, newStatus return_request_status.Status) error
	GetPendingReturnRequests(userID int64) ([]models.ReturnRequest, error)
}
