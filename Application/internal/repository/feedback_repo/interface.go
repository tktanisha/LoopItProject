//go:generate mockgen -source=interface.go -destination=../../mock/mock_feedback_repo.go -package=mock
package feedback_repo

import "loopit/internal/models"

type FeedbackRepository interface {
	CreateFeedback(feedback models.Feedback) error
	GetAllFeedbacks() ([]models.Feedback, error)
	Save() error
}
