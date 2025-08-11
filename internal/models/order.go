package models

import (
	"loopit/internal/enums/order_status"
	"time"
)

type Order struct {
	ID             int                 `json:"id"`
	ProductID      int                 `json:"product_id"`
	UserID         int                 `json:"user_id"`
	StartDate      time.Time           `json:"start_date"`
	EndDate        time.Time           `json:"end_date"`
	TotalAmount    float64             `json:"total_amount"`
	SecurityAmount float64             `json:"security_amount"`
	Status         order_status.Status `json:"status"`
	CreatedAt      time.Time           `json:"created_at"`
}
