package models

import (
	"loopit/internal/enums/buyer_request_status"
	"time"
)

type BuyingRequest struct {
	ID          int64                         `json:"id"`
	ProductID   int64                        `json:"product_id"`
	RequestedBy int64                         `json:"requested_by"`
	Status      buyer_request_status.Status `json:"status"`
	CreatedAt   time.Time                   `json:"created_at"`
}

type BuyingRequestDto struct {
	BuyRequest BuyingRequest   `json:"buy_request"`
	Product    ProductResponse `json:"product"`
}
