package models

type Lender struct {
	ID            int64     `json:"id"`
	IsVerified    bool    `json:"is_verified"`
	TotalEarnings float64 `json:"total_earnings"`
}
