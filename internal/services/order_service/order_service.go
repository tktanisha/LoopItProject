package order_service

import (
	"errors"
	"fmt"
	"loopit/internal/enums"
	"loopit/internal/enums/order_status"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
	"loopit/pkg/logger"
)

type OrderService struct {
	orderRepo         order_repo.OrderRepo
	returnRequestRepo return_request_repo.ReturnRequestRepo
	productRepo       product_repo.ProductRepo
	log               *logger.Logger
}

func NewOrderService(
	OrderRepo order_repo.OrderRepo,
	returnRepo return_request_repo.ReturnRequestRepo,
	productRepo product_repo.ProductRepo,
	log *logger.Logger,
) OrderServiceInterface {
	return &OrderService{
		orderRepo:         OrderRepo,
		returnRequestRepo: returnRepo,
		productRepo:       productRepo,
		log:               log,
	}
}

// in use(on create order), return-requested(on create return request), returned(on mark as returned)
func (s *OrderService) UpdateOrderStatus(orderID int, newStatus order_status.Status) error {
	s.log.Info(fmt.Sprintf("Updating status of order %d to %s", orderID, newStatus))
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch order %d, error: %v", orderID, err))
		return err
	}
	if order == nil {
		s.log.Warning(fmt.Sprintf("Order %d not found", orderID))
		return errors.New("order not found")
	}

	if newStatus == order_status.Returned && order.Status != order_status.ReturnRequested {
		s.log.Warning(fmt.Sprintf("Invalid transition: order %d is in %s, cannot mark as returned", orderID, order.Status))
		return errors.New("order must be in return_requested status to mark as returned")
	}

	if err := s.orderRepo.UpdateOrderStatus(orderID, newStatus.String()); err != nil {
		s.log.Error(fmt.Sprintf("Failed to update order %d status to %s, error: %v", orderID, newStatus, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Successfully updated order %d status to %s", orderID, newStatus))
	return nil
}

// GetOrderHistory returns orders for the user based on role
func (s *OrderService) GetOrderHistory(userCtx *models.UserContext, filterStatus []order_status.Status) ([]*models.Order, error) {
	s.log.Info(fmt.Sprintf("Fetching order history for user %d with role %s", userCtx.ID, userCtx.Role))
	filterStatusStr := []string{}
	for _, status := range filterStatus {
		filterStatusStr = append(filterStatusStr, status.String())
	}

	orders, err := s.orderRepo.GetOrderHistory(userCtx.ID, filterStatusStr)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch order history for user %d, error: %v", userCtx.ID, err))
		return nil, err
	}
	s.log.Info(fmt.Sprintf("Fetched %d orders for user %d", len(orders), userCtx.ID))
	return orders, nil
}

func (s *OrderService) GetLenderOrders(userCtx *models.UserContext) ([]*models.Order, error) {
	s.log.Info(fmt.Sprintf("Fetching lender orders for user %d", userCtx.ID))
	if userCtx.Role != enums.RoleLender {
		s.log.Warning(fmt.Sprintf("Unauthorized attempt: user %d with role %s tried to fetch lender orders", userCtx.ID, userCtx.Role))
		return nil, errors.New("only lender can get orders")
	}

	orders, err := s.orderRepo.GetLenderOrders(userCtx.ID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch lender orders for user %d, error: %v", userCtx.ID, err))
		return nil, err
	}
	s.log.Info(fmt.Sprintf("Fetched %d lender orders for user %d", len(orders), userCtx.ID))
	return orders, nil
}

// Lender marks order as returned (status: returned)
func (s *OrderService) MarkOrderAsReturned(orderID int, userCtx *models.UserContext) error {
	s.log.Info(fmt.Sprintf("User %d attempting to mark order %d as returned", userCtx.ID, orderID))
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch order %d, error: %v", orderID, err))
		return err
	}
	if order == nil {
		s.log.Warning(fmt.Sprintf("Order %d not found", orderID))
		return errors.New("order not found")
	}

	product, err := s.productRepo.FindByID(order.ProductID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch product %d for order %d, error: %v", order.ProductID, orderID, err))
		return errors.New("unable to find product for the order")
	}
	if product == nil {
		s.log.Warning(fmt.Sprintf("Product %d for order %d not found", order.ProductID, orderID))
		return errors.New("product not found")
	}
	if product.Product.LenderID != userCtx.ID {
		s.log.Warning(fmt.Sprintf("Unauthorized attempt: user %d tried to mark order %d as returned (not lender)", userCtx.ID, orderID))
		return errors.New("unauthorized lender")
	}

	returnRequests, err := s.returnRequestRepo.GetAllReturnRequests([]string{return_request_status.Approved.String()})
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch return requests for order %d, error: %v", orderID, err))
		return errors.New("unable to find return requests for the order")
	}

	isExists := false
	for _, rr := range returnRequests {
		if rr.OrderID == orderID {
			isExists = true
			break
		}
	}
	if !isExists {
		s.log.Warning(fmt.Sprintf("No approved return request found for order %d", orderID))
		return errors.New("order has not been approved for return")
	}

	if err := s.orderRepo.UpdateOrderStatus(orderID, order_status.Returned.String()); err != nil {
		s.log.Error(fmt.Sprintf("Failed to update order %d status to returned, error: %v", orderID, err))
		return err
	}
	s.log.Info(fmt.Sprintf("Order %d successfully marked as returned by lender %d", orderID, userCtx.ID))
	return nil
}

// GetAllApprovedAwaitingOrders returns orders approved for return but not yet marked returned
func (s *OrderService) GetAllApprovedAwaitingOrders(userCtx *models.UserContext) ([]*models.Order, error) {
	s.log.Info(fmt.Sprintf("Fetching approved awaiting orders for lender %d", userCtx.ID))
	if userCtx.Role != enums.RoleLender {
		s.log.Warning(fmt.Sprintf("Unauthorized attempt: user %d with role %s tried to fetch awaiting orders", userCtx.ID, userCtx.Role))
		return nil, errors.New("only lender can get returned awaiting orders")
	}

	returnRequests, err := s.returnRequestRepo.GetAllReturnRequests([]string{return_request_status.Approved.String()})
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch return requests, error: %v", err))
		return nil, errors.New("unable to find return requests for the order")
	}

	orders := []*models.Order{}
	for _, rr := range returnRequests {
		order, err := s.orderRepo.GetOrderByID(rr.OrderID)
		if err != nil {
			s.log.Error(fmt.Sprintf("Failed to fetch order %d for return request, error: %v", rr.OrderID, err))
			return nil, errors.New("unable to find order for the return request")
		}
		if order == nil {
			s.log.Warning(fmt.Sprintf("Order %d referenced in return request not found", rr.OrderID))
			continue
		}
		orders = append(orders, order)
	}

	s.log.Info(fmt.Sprintf("Fetched %d approved awaiting orders for lender %d", len(orders), userCtx.ID))
	return orders, nil
}
