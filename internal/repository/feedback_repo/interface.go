package feedback_repo

import "loopit/internal/models"

type FeedbackRepository interface {
	CreateFeedback(feedback models.Feedback) error
	GetAllFeedbacks() ([]models.Feedback, error)
	Save() error
}
