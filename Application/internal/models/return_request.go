package models

import (
	"loopit/internal/enums/return_request_status"
	"time"
)

type ReturnRequest struct {
	ID          int64                          `json:"id"`
	OrderID     int64                          `json:"order_id"`
	RequestedBy int64                          `json:"requested_by"`
	Status      return_request_status.Status `json:"status"`
	CreatedAt   time.Time                    `json:"created_at"`
}


type ReturnRequestDto struct{
	ReturnRequest ReturnRequest  `json:"buy_request"`
	Product    ProductResponse `json:"product"`
}