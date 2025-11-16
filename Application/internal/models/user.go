package models

import (
	"loopit/internal/enums"
	"time"
)


type User struct {
	ID           int64       `json:"id" dynamodbav:"UserID"`
	FullName     string     `json:"full_name" dynamodbav:"FullName"`
	Email        string     `json:"email" dynamodbav:"Email"`
	PhoneNumber  string     `json:"phone_number" dynamodbav:"PhoneNumber"`
	Address      string     `json:"address" dynamodbav:"Address"`
	PasswordHash string     `json:"password_hash" dynamodbav:"PasswordHash"`
	SocietyID    int64        `json:"society_id" dynamodbav:"SocietyID"`
	Role         enums.Role `json:"role" dynamodbav:"Role"`
	CreatedAt    time.Time    `json:"created_at" dynamodbav:"CreatedAt"`
	PK           string     `dynamodbav:"pk"`
	SK           string     `dynamodbav:"sk"`
}

type UserFilter struct {
	Search    string  `json:"search"`
	Role      string   `json:"role"`
	SocietyID string   `json:"society_id"`
}

type UserContext struct {
	ID   int64
	Name string
	Role enums.Role
}
