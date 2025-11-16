package models

import "time"

type ProductImage struct {
	ID         int64       `json:"id"`
	ProductID  int64       `json:"product_id"`
	ImageURL   string    `json:"image_url"`
	UploadedAt time.Time `json:"uploaded_at"`
}
