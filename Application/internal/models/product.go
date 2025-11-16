package models

import "time"

type Product struct {
	ID          int64       `json:"id"`
	LenderID    int64       `json:"lender_id"`
	CategoryID  int64       `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"` // Can be time.Duration or string
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	ImageUrl   string     `json:"image_url"`
}

type ProductResponse struct {
	Product  Product  `json:"product"`
	Category Category `json:"category"`
	User     User     `json:"user"`
}

type ProductFilter struct {
	Search      string
	LenderID    string
	CategoryID  string
	IsAvailable string
}