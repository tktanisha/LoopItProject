package models

import (
	"loopit/internal/enums/buyer_request_status"
	"time"
)

type BuyingRequest struct {
	ID          int                         `json:"id"`
	ProductID   int                         `json:"product_id"`
	RequestedBy int                         `json:"requested_by"`
	Status      buyer_request_status.Status `json:"status"`
	CreatedAt   time.Time                   `json:"created_at"`
}
