package models

import "time"

type Product struct {
	ID          int       `json:"id"`
	LenderID    int       `json:"lender_id"`
	CategoryID  int       `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"` // Can be time.Duration or string
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductResponse struct {
	Product  Product  `json:"product"`
	Category Category `json:"category"`
	User     User     `json:"user"`
}
