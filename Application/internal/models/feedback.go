package models

import "time"

type Feedback struct {
	ID        int64       `json:"id"`
	GivenBy   int64       `json:"given_by"`
	GivenTo   int64      `json:"given_to"`
	Text      string    `json:"text"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}
