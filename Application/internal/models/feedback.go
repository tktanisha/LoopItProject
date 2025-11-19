package models

import "time"

type Feedback struct {
	ID        int64     `json:"id,string" dynamodbav:"ID"`
	GivenBy   int64     `json:"given_by,string" dynamodbav:"GivenBy"`
	GivenTo   int64     `json:"given_to,string" dynamodbav:"GivenTo"`
	Text      string    `json:"text" dynamodbav:"Text"`
	Rating    int       `json:"rating" dynamodbav:"Rating"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"CreatedAt"`
}
