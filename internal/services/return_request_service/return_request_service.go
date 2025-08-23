package return_request_service

import (
	"errors"
	"fmt"
	"loopit/internal/enums/order_status"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
	"loopit/pkg/logger"
	"time"
)

type returnRequestService struct {
	orderRepo         order_repo.OrderRepo
	productRepo       product_repo.ProductRepo
	returnRequestRepo return_request_repo.ReturnRequestRepo
	log               *logger.Logger
}

func NewReturnRequestService(
	orderRepo order_repo.OrderRepo,
	productRepo product_repo.ProductRepo,
	rrRepo return_request_repo.ReturnRequestRepo,
	log *logger.Logger,
) ReturnRequestServiceInterface {
	return &returnRequestService{
		orderRepo:         orderRepo,
		returnRequestRepo: rrRepo,
		productRepo:       productRepo,
		log:               log,
	}
}

func (s *returnRequestService) CreateReturnRequest(userID int, orderID int) error {
	s.log.Info(fmt.Sprintf("User %d attempting to create return request for order %d", userID, orderID))

	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch order %d: %v", orderID, err))
		return err
	}

	if order.Status != order_status.InUse {
		s.log.Warning(fmt.Sprintf("Order %d is not in 'in_use' status", orderID))
		return errors.New("order is not in 'in_use' status")
	}

	productID := order.ProductID
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch product %d: %v", productID, err))
		return err
	}

	if product.Product.LenderID != userID {
		s.log.Warning(fmt.Sprintf("User %d is not lender of product %d", userID, productID))
		return errors.New("user is not lender of the order's product")
	}

	returnRequest := models.ReturnRequest{
		OrderID:     orderID,
		RequestedBy: userID,
		Status:      return_request_status.Pending,
		CreatedAt:   time.Now(),
	}

	if err := s.orderRepo.UpdateOrderStatus(orderID, order_status.ReturnRequested.String()); err != nil {
		s.log.Error(fmt.Sprintf("Failed to update order %d status: %v", orderID, err))
		return err
	}

	if err := s.returnRequestRepo.CreateReturnRequest(returnRequest); err != nil {
		s.log.Error(fmt.Sprintf("Failed to create return request for order %d: %v", orderID, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Return request created successfully for order %d by user %d", orderID, userID))
	return nil
}

func (s *returnRequestService) UpdateReturnRequestStatus(userID int, reqID int, newStatus return_request_status.Status) error {
	s.log.Info(fmt.Sprintf("User %d attempting to update return request %d to status %s", userID, reqID, newStatus.String()))

	if newStatus != return_request_status.Approved && newStatus != return_request_status.Rejected {
		s.log.Warning(fmt.Sprintf("Invalid status update attempt by user %d for request %d: %s", userID, reqID, newStatus.String()))
		return errors.New("invalid status update")
	}

	req, err := s.returnRequestRepo.GetReturnRequestByID(reqID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch return request %d: %v", reqID, err))
		return err
	}

	if req.Status != return_request_status.Pending {
		s.log.Warning(fmt.Sprintf("Return request %d not in pending status", reqID))
		return errors.New("return request is not in pending status")
	}

	order, err := s.orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch order %d for return request %d: %v", req.OrderID, reqID, err))
		return err
	}

	if order.UserID != userID {
		s.log.Warning(fmt.Sprintf("User %d does not own order %d linked to return request %d", userID, order.ID, reqID))
		return errors.New("user does not own this order")
	}

	if err := s.returnRequestRepo.UpdateReturnRequestStatus(req.ID, newStatus.String()); err != nil {
		s.log.Error(fmt.Sprintf("Failed to update return request %d status: %v", reqID, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Return request %d updated to status %s by user %d", reqID, newStatus.String(), userID))
	return nil
}

func (s *returnRequestService) GetPendingReturnRequests(userID int) ([]models.ReturnRequest, error) {
	s.log.Info(fmt.Sprintf("Fetching pending return requests for user %d", userID))

	allRequests, err := s.returnRequestRepo.GetAllReturnRequests([]string{return_request_status.Pending.String()})
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch return requests: %v", err))
		return nil, err
	}

	var userRequests []models.ReturnRequest
	for _, req := range allRequests {
		order, err := s.orderRepo.GetOrderByID(req.OrderID)
		if err != nil {
			s.log.Error(fmt.Sprintf("Failed to fetch order %d for return request %d: %v", req.OrderID, req.ID, err))
			return nil, err
		}
		if order.UserID == userID {
			userRequests = append(userRequests, req)
		}
	}

	s.log.Info(fmt.Sprintf("User %d has %d pending return requests", userID, len(userRequests)))
	return userRequests, nil
}
