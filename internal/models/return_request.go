package models

import (
	"loopit/internal/enums/return_request_status"
	"time"
)

type ReturnRequest struct {
	ID          int                          `json:"id"`
	OrderID     int                          `json:"order_id"`
	RequestedBy int                          `json:"requested_by"`
	Status      return_request_status.Status `json:"status"`
	CreatedAt   time.Time                    `json:"created_at"`
}
