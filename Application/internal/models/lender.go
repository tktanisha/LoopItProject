package models

type Lender struct {
	ID            int64   `json:"id,string" dynamodbav:"ID"`
	IsVerified    bool    `json:"is_verified" dynamodbav:"IsVerified"`
	TotalEarnings float64 `json:"total_earnings" dynamodbav:"TotalEarnings"`
}
