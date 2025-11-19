package models

type Category struct {
	ID       int64   `json:"id,string" dynamodbav:"ID"`
	Name     string  `json:"name" dynamodbav:"Name"`
	Price    float64 `json:"price" dynamodbav:"Price"`
	Security float64 `json:"security" dynamodbav:"Security"`
}
