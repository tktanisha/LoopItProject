package models

import "time"

type Feedback struct {
	ID        int       `json:"id"`
	GivenBy   int       `json:"given_by"`
	GivenTo   int       `json:"given_to"`
	Text      string    `json:"text"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}
