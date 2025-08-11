package models

import "time"

type ProductImage struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"product_id"`
	ImageURL   string    `json:"image_url"`
	UploadedAt time.Time `json:"uploaded_at"`
}
