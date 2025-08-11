package models

import (
	"loopit/internal/enums"
	"time"
)

type User struct {
	ID           int        `json:"id"`
	FullName     string     `json:"full_name"`
	Email        string     `json:"email"`
	PhoneNumber  string     `json:"phone_number"`
	Address      string     `json:"address"`
	PasswordHash string     `json:"password_hash"`
	SocietyID    int        `json:"society_id"`
	Role         enums.Role `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
}

type UserContext struct {
	ID   int
	Name string
	Role enums.Role
}
