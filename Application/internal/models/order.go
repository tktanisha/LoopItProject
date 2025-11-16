package models

import (
	"loopit/internal/enums/order_status"
	"time"
)

type Order struct {
	ID             int64               `json:"id" dynamodbav:"ID"`
	ProductID      int64               `json:"product_id" dynamodbav:"ProductID"`
	UserID         int64               `json:"user_id" dynamodbav:"UserID"`
	StartDate      time.Time           `json:"start_date" dynamodbav:"StartDate"`
	EndDate        time.Time           `json:"end_date" dynamodbav:"EndDate"`
	TotalAmount    float64             `json:"total_amount" dynamodbav:"TotalAmount"`
	SecurityAmount float64             `json:"security_amount" dynamodbav:"SecurityAmount"`
	Status         order_status.Status `json:"status" dynamodbav:"Status"`
	CreatedAt      time.Time           `json:"created_at" dynamodbav:"CreatedAt"`
}

type OrderDto struct {
	Order   Order           `json:"order"`
	Product ProductResponse `json:"product"`
}
