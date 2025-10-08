package buyer_request_service

import (
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
)

type BuyerRequestServiceInterface interface {
	CreateBuyerRequest(productID int, userCtx *models.UserContext) error
	UpdateBuyerRequestStatus(requestID int, updatedStatus buyer_request_status.Status, userCtx *models.UserContext) error
	GetAllBuyerRequests(productID *int, status []string) ([]models.BuyingRequest, error)
}
