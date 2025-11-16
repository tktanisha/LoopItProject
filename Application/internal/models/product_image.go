package models

import "time"

type ProductImage struct {
	ID         int64     `json:"id" dynamodbav:"ID"`
	ProductID  int64     `json:"product_id" dynamodbav:"ProductID"`
	ImageURL   string    `json:"image_url" dynamodbav:"ImageURL"`
	UploadedAt time.Time `json:"uploaded_at" dynamodbav:"UploadedAt"`
}
