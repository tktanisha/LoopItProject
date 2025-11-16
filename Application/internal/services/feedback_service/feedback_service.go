package feedback_service

import (
	"errors"
	"fmt"
	"log"
	"loopit/internal/models"
	"loopit/internal/repository/feedback_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"time"
)

type FeedbackService struct {
	feedback_repo feedback_repo.FeedbackRepository
	product_repo  product_repo.ProductRepo
	order_repo    order_repo.OrderRepo
}

func NewFeedbackService(
	repo feedback_repo.FeedbackRepository,
	productRepo product_repo.ProductRepo,
	orderRepo order_repo.OrderRepo,
) FeedbackServiceInterface {
	return &FeedbackService{
		feedback_repo: repo,
		product_repo:  productRepo,
		order_repo:    orderRepo,
	
	}
}

func (s *FeedbackService) GiveFeedback(orderID int64, feedbackText string, rating int, userCtx *models.UserContext) error {
	log.Print(fmt.Sprintf("User %d attempting to give feedback for order %d", userCtx.ID, orderID))

	order, err := s.order_repo.GetOrderByID(orderID)
	if err != nil {
		log.Print(fmt.Sprintf("Order not found for ID %d, error: %v", orderID, err))
		return err
	}
	product, err := s.product_repo.FindByID(order.ProductID)
	if err != nil {
		log.Print(fmt.Sprintf("Product not found for order %d, productID %d, error: %v", orderID, order.ProductID, err))
		return err
	}

	givenTo := product.Product.LenderID
	if userCtx.ID == givenTo {
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
		return err
	}

	return nil
}

func (s *FeedbackService) GetAllGivenFeedbacks(userCtx *models.UserContext) ([]models.Feedback, error) {

	feedbacks, err := s.feedback_repo.GetAllFeedbacks()
	if err != nil {
		return nil, err
	}

	var givenFeedbacks []models.Feedback
	for _, feedback := range feedbacks {
		if feedback.GivenBy == userCtx.ID {
			givenFeedbacks = append(givenFeedbacks, feedback)
		}
	}

	return givenFeedbacks, nil
}

func (s *FeedbackService) GetAllReceivedFeedbacks(userCtx *models.UserContext) ([]models.Feedback, error) {

	feedbacks, err := s.feedback_repo.GetAllFeedbacks()
	if err != nil {
		return nil, err
	}

	var receivedFeedbacks []models.Feedback
	for _, feedback := range feedbacks {
		if feedback.GivenTo == userCtx.ID {
			receivedFeedbacks = append(receivedFeedbacks, feedback)
		}
	}

	return receivedFeedbacks, nil
}
