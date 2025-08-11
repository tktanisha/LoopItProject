package models

import "time"

type Society struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Pincode   string    `json:"pincode"`
	CreatedAt time.Time `json:"created_at"`
}
