package models

import (
	"loopit/internal/enums/order_status"
	"time"
)

type Order struct {
	ID             int64                 `json:"id"`
	ProductID      int64                 `json:"product_id"`
	UserID         int64                 `json:"user_id"`
	StartDate      time.Time           `json:"start_date"`
	EndDate        time.Time           `json:"end_date"`
	TotalAmount    float64             `json:"total_amount"`
	SecurityAmount float64             `json:"security_amount"`
	Status         order_status.Status `json:"status"`
	CreatedAt      time.Time           `json:"created_at"`
}

type OrderDto struct {
	Order   Order           `json:"order"`
	Product ProductResponse `json:"product"`
}
