
package models

import "time"

type Society struct {
    ID        int64     `json:"id" dynamodbav:"ID"`
    Name      string    `json:"name" dynamodbav:"Name"`
    Location  string    `json:"location" dynamodbav:"Location"`
    Pincode   string    `json:"pincode" dynamodbav:"Pincode"`
    CreatedAt time.Time `json:"created_at" dynamodbav:"CreatedAt"`
}