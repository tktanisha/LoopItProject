package feedback_repo

import (
	"database/sql"
	"fmt"
	"loopit/internal/models"
	"loopit/pkg/logger"
	"time"
)

type FeedBackDBRepo struct {
	db  *sql.DB
	log *logger.Logger
}

func NewFeedBackDBRepo(db *sql.DB, log *logger.Logger) *FeedBackDBRepo {
	return &FeedBackDBRepo{db: db, log: log}
}

// CreateFeedback inserts a new feedback into the database
func (r *FeedBackDBRepo) CreateFeedback(feedback models.Feedback) error {
	query := `
    INSERT INTO feedbacks (given_by, given_to, text, rating, created_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id
    `
	err := r.db.QueryRow(query, feedback.GivenBy, feedback.GivenTo, feedback.Text, feedback.Rating, time.Now()).Scan(&feedback.ID)
	if err != nil && r.log != nil {
		r.log.Error(fmt.Sprintf("Repo: DB error creating feedback by user %d for user %d: %v", feedback.GivenBy, feedback.GivenTo, err))
	}
	return err
}

// GetAllFeedbacks returns all feedbacks from the database
func (r *FeedBackDBRepo) GetAllFeedbacks() ([]models.Feedback, error) {
	rows, err := r.db.Query("SELECT id, given_by, given_to, text, rating, created_at FROM feedbacks")
	if err != nil {
		if r.log != nil {
			r.log.Error(fmt.Sprintf("Repo: DB error fetching all feedbacks: %v", err))
		}
		return nil, err
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var f models.Feedback
		if err := rows.Scan(&f.ID, &f.GivenBy, &f.GivenTo, &f.Text, &f.Rating, &f.CreatedAt); err != nil {
			if r.log != nil {
				r.log.Warning(fmt.Sprintf("Repo: Could not scan feedback row: %v", err))
			}
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
