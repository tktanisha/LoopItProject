package models

import (
	"loopit/internal/enums/buyer_request_status"
	"time"
)

type BuyingRequest struct {
	ID          int64                       `json:"id,string" dynamodbav:"ID"`
	ProductID   int64                       `json:"product_id,string" dynamodbav:"ProductId"`
	RequestedBy int64                       `json:"requested_by,string" dynamodbav:"RequestedBy"`
	Status      buyer_request_status.Status `json:"status" dynamodbav:"Status"`
	CreatedAt   time.Time                   `json:"created_at" dynamodbav:"CreatedAt"`
}

type BuyingRequestDto struct {
	BuyRequest BuyingRequest   `json:"buy_request"`
	Product    ProductResponse `json:"product"`
}
