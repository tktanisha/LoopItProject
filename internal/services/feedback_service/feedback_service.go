package feedback_service

import (
	"errors"
	"fmt"
	"loopit/internal/enums/order_status"
	"loopit/internal/models"
	"loopit/internal/repository/feedback_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/pkg/logger"
	"time"
)

type FeedbackService struct {
	feedback_repo feedback_repo.FeedbackRepository
	product_repo  product_repo.ProductRepo
	order_repo    order_repo.OrderRepo
	log           *logger.Logger
}

func NewFeedbackService(
	repo feedback_repo.FeedbackRepository,
	productRepo product_repo.ProductRepo,
	orderRepo order_repo.OrderRepo,
	log *logger.Logger,
) FeedbackServiceInterface {
	return &FeedbackService{
		feedback_repo: repo,
		product_repo:  productRepo,
		order_repo:    orderRepo,
		log:           log,
	}
}

func (s *FeedbackService) GiveFeedback(orderID int, feedbackText string, rating int, userCtx *models.UserContext) error {
	s.log.Info(fmt.Sprintf("User %d attempting to give feedback for order %d", userCtx.ID, orderID))

	order, err := s.order_repo.GetOrderByID(orderID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Order not found for ID %d, error: %v", orderID, err))
		return err
	}

	if order.Status != order_status.Returned {
		s.log.Warning(fmt.Sprintf("Feedback rejected for order %d: order not returned (status: %s)", orderID, order.Status))
		return errors.New("feedback can only be given for returned orders")
	}

	product, err := s.product_repo.FindByID(order.ProductID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Product not found for order %d, productID %d, error: %v", orderID, order.ProductID, err))
		return err
	}

	givenTo := product.Product.LenderID
	if userCtx.ID == givenTo {
		s.log.Warning(fmt.Sprintf("User %d attempted to give feedback to self for order %d", userCtx.ID, orderID))
		return errors.New("you cannot give feedback to yourself")
	}

	feedback := models.Feedback{
		GivenBy:   userCtx.ID,
		GivenTo:   givenTo,
		Rating:    rating,
		Text:      feedbackText,
		CreatedAt: time.Now(),
	}

	if err := s.feedback_repo.CreateFeedback(feedback); err != nil {
		s.log.Error(fmt.Sprintf("Failed to save feedback for order %d by user %d, error: %v", orderID, userCtx.ID, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Feedback successfully given by user %d to user %d for order %d", userCtx.ID, givenTo, orderID))
	return nil
}

func (s *FeedbackService) GetAllGivenFeedbacks(userCtx *models.UserContext) ([]models.Feedback, error) {
	s.log.Info(fmt.Sprintf("Fetching all feedbacks given by user %d", userCtx.ID))

	feedbacks, err := s.feedback_repo.GetAllFeedbacks()
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch feedbacks, error: %v", err))
		return nil, err
	}

	var givenFeedbacks []models.Feedback
	for _, feedback := range feedbacks {
		if feedback.GivenBy == userCtx.ID {
			givenFeedbacks = append(givenFeedbacks, feedback)
		}
	}

	s.log.Info(fmt.Sprintf("User %d has given %d feedback(s)", userCtx.ID, len(givenFeedbacks)))
	return givenFeedbacks, nil
}

func (s *FeedbackService) GetAllReceivedFeedbacks(userCtx *models.UserContext) ([]models.Feedback, error) {
	s.log.Info(fmt.Sprintf("Fetching all feedbacks received by user %d", userCtx.ID))

	feedbacks, err := s.feedback_repo.GetAllFeedbacks()
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch feedbacks, error: %v", err))
		return nil, err
	}

	var receivedFeedbacks []models.Feedback
	for _, feedback := range feedbacks {
		if feedback.GivenTo == userCtx.ID {
			receivedFeedbacks = append(receivedFeedbacks, feedback)
		}
	}

	s.log.Info(fmt.Sprintf("User %d has received %d feedback(s)", userCtx.ID, len(receivedFeedbacks)))
	return receivedFeedbacks, nil
}
