package models

import "time"

type Product struct {
	ID          int64     `json:"id,string" dynamodbav:"ID"`
	LenderID    int64     `json:"lender_id,string" dynamodbav:"LenderID"`
	CategoryID  int64     `json:"category_id,string" dynamodbav:"CategoryID"`
	Name        string    `json:"name" dynamodbav:"Name"`
	Description string    `json:"description" dynamodbav:"Description"`
	Duration    int       `json:"duration" dynamodbav:"Duration"`
	IsAvailable bool      `json:"is_available" dynamodbav:"IsAvailable"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"CreatedAt"`
	ImageUrl    string    `json:"image_url" dynamodbav:"ImageUrl"`
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
