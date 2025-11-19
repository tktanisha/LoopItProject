package models

import (
	"loopit/internal/enums/return_request_status"
	"time"
)

type ReturnRequest struct {
	ID          int64                        `json:"id,string" dynamodbav:"ID"`
	OrderID     int64                        `json:"order_id,string" dynamodbav:"OrderID"`
	RequestedBy int64                        `json:"requested_by,string" dynamodbav:"RequestedBy"`
	Status      return_request_status.Status `json:"status" dynamodbav:"Status"`
	CreatedAt   time.Time                    `json:"created_at" dynamodbav:"CreatedAt"`
}

type ReturnRequestDto struct {
	ReturnRequest ReturnRequest   `json:"buy_request"`
	Product       ProductResponse `json:"product"`
}
