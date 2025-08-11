package return_request_service

import (
	"errors"
	"loopit/internal/enums/order_status"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
	"time"
)

type returnRequestService struct {
	orderRepo         order_repo.OrderRepo
	productRepo       product_repo.ProductRepo
	returnRequestRepo return_request_repo.ReturnRequestRepo
}

func NewReturnRequestService(orderRepo order_repo.OrderRepo, productRepo product_repo.ProductRepo, rrRepo return_request_repo.ReturnRequestRepo) ReturnRequestServiceInterface {
	return &returnRequestService{
		orderRepo:         orderRepo,
		returnRequestRepo: rrRepo,
		productRepo:       productRepo,
	}
}

func (s *returnRequestService) CreateReturnRequest(userID int, orderID int) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	// Must be "in_use"
	if order.Status != order_status.InUse {
		return errors.New("order is not in 'in_use' status")
	}

	productID := order.ProductID
	product, err := s.productRepo.FindByID(productID)

	if err != nil {
		return err
	}

	if product.Product.LenderID != userID {
		return errors.New("user is not lender of the order's product")
	}

	if order.Status != order_status.InUse {
		return errors.New("order is not in 'in_use' status and cannot be returned")
	}

	returnRequest := models.ReturnRequest{
		OrderID:     orderID,
		RequestedBy: userID,
		Status:      return_request_status.Pending,
		CreatedAt:   time.Now(),
	}

	if err := s.orderRepo.UpdateOrderStatus(orderID, order_status.ReturnRequested.String()); err != nil {
		return err
	}

	return s.returnRequestRepo.CreateReturnRequest(returnRequest)
}

func (s *returnRequestService) UpdateReturnRequestStatus(userID int, reqID int, newStatus return_request_status.Status) error {
	// Only accept or reject allowed
	if newStatus != return_request_status.Approved && newStatus != return_request_status.Rejected {
		return errors.New("invalid status update")
	}

	req, err := s.returnRequestRepo.GetReturnRequestByID(reqID)
	if err != nil {
		return err
	}

	if req.Status != return_request_status.Pending {
		return errors.New("return request is not in pending status")
	}

	order, err := s.orderRepo.GetOrderByID(req.OrderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("user does not own this order")
	}
	return s.returnRequestRepo.UpdateReturnRequestStatus(req.ID, newStatus.String())
}

func (s *returnRequestService) GetPendingReturnRequests(userID int) ([]models.ReturnRequest, error) {
	allRequests, _ := s.returnRequestRepo.GetAllReturnRequests([]string{return_request_status.Pending.String()})
	var userRequests []models.ReturnRequest
	for _, req := range allRequests {
		orderId := req.OrderID
		order, err := s.orderRepo.GetOrderByID(orderId)
		if err != nil {
			return nil, err
		}
		if order.UserID == userID {
			userRequests = append(userRequests, req)
		}
	}
	return userRequests, nil
}
