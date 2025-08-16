package feedback_repo

import (
	"database/sql"
	"loopit/internal/models"
	"time"
)

type FeedBackDBRepo struct {
	db *sql.DB
}

func NewFeedBackDBRepo(db *sql.DB) *FeedBackDBRepo {
	return &FeedBackDBRepo{db: db}
}

// CreateFeedback inserts a new feedback into the database
func (r *FeedBackDBRepo) CreateFeedback(feedback models.Feedback) error {
	query := `
	INSERT INTO feedbacks (given_by, given_to, text, rating, created_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	return r.db.QueryRow(query, feedback.GivenBy, feedback.GivenTo, feedback.Text, feedback.Rating, time.Now()).Scan(feedback.ID)
}

// GetAllFeedbacks returns all feedbacks from the database
func (r *FeedBackDBRepo) GetAllFeedbacks() ([]models.Feedback, error) {
	rows, err := r.db.Query("SELECT id, given_by, given_to, text, rating, created_at FROM feedbacks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var f models.Feedback
		if err := rows.Scan(&f.ID, &f.GivenBy, &f.GivenTo, &f.Text, &f.Rating, &f.CreatedAt); err != nil {
			continue
		}
		feedbacks = append(feedbacks, f)
	}

	if len(feedbacks) == 0 {
		return nil, nil
	}

	return feedbacks, nil
}

// Save is a no-op for Postgres
func (r *FeedBackDBRepo) Save() error {
	return nil
}
