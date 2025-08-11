package feedback_service

import "loopit/internal/models"

type FeedbackServiceInterface interface {
	GiveFeedback(orderID int, feedbackText string, rating int, userCtx *models.UserContext) error
	GetAllGivenFeedbacks(userCtx *models.UserContext) ([]models.Feedback, error)
	GetAllReceivedFeedbacks(userCtx *models.UserContext) ([]models.Feedback, error)
}
