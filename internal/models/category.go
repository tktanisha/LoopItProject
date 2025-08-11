package models

type Category struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Security float64 `json:"security"`
}
