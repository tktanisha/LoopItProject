package order_service

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/enums/order_status"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
)

type OrderService struct {
	orderRepo         order_repo.OrderRepo
	returnRequestRepo return_request_repo.ReturnRequestRepo
	productRepo       product_repo.ProductRepo
}

func NewOrderService(OrderRepo order_repo.OrderRepo, returnRepo return_request_repo.ReturnRequestRepo, productRepo product_repo.ProductRepo) OrderServiceInterface {
	return &OrderService{
		orderRepo:         OrderRepo,
		returnRequestRepo: returnRepo,
		productRepo:       productRepo,
	}
}

// in use(on create order),return-requested(on create return request),returned(on mark as returned)
func (r *OrderService) UpdateOrderStatus(orderID int, newStatus order_status.Status) error {
	order, err := r.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	// TODO try to remove this check to reduce tight coupling
	if newStatus == order_status.Returned && order.Status != order_status.ReturnRequested {
		return errors.New("order must be in return_requested status to mark as returned")
	}

	return r.orderRepo.UpdateOrderStatus(orderID, newStatus.String())

}

//	TODO:  Get Order History -> user, lender
//
// GetOrderHistory returns orders for the user based on role
// if user is lender then return all orders where user is lender
// if user is buyer then return all orders where user is buyer
func (s *OrderService) GetOrderHistory(userCtx *models.UserContext, filterStatus []order_status.Status) ([]*models.Order, error) {
	filterStatusStr := []string{}
	for _, status := range filterStatus {
		filterStatusStr = append(filterStatusStr, status.String())
	}
	return s.orderRepo.GetOrderHistory(userCtx.ID, filterStatusStr)
}

func (s *OrderService) GetLenderOrders(userCtx *models.UserContext) ([]*models.Order, error) {
	if userCtx.Role != enums.RoleLender {
		return nil, errors.New("only lender can get orders")
	}

	return s.orderRepo.GetLenderOrders(userCtx.ID)
}

// Lender marks order as returned (status: returned)
func (s *OrderService) MarkOrderAsReturned(orderID int, userCtx *models.UserContext) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("order not found")
	}

	product, err := s.productRepo.FindByID(order.ProductID)
	if err != nil {
		return errors.New("unable to find product for the order")
	}
	if product == nil {
		return errors.New("product not found")
	}
	if product.Product.LenderID != userCtx.ID {
		return errors.New("unauthorized lender")
	}

	returnRequests, err := s.returnRequestRepo.GetAllReturnRequests([]string{return_request_status.Approved.String()})
	if err != nil {
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
		return errors.New("order has not been approved for return")
	}

	return s.orderRepo.UpdateOrderStatus(orderID, order_status.Returned.String())
}

// getAllConfirmedAwaitingOrders returns orders that are approved by user for return requested and waiting for return status in order
func (s *OrderService) GetAllApprovedAwaitingOrders(userCtx *models.UserContext) ([]*models.Order, error) {
	if userCtx.Role != enums.RoleLender {
		return nil, errors.New("only lender can get returned awaiting orders")
	}

	returnRequests, err := s.returnRequestRepo.GetAllReturnRequests([]string{return_request_status.Approved.String()})
	if err != nil {
		return nil, errors.New("unable to find return requests for the order")
	}

	orders := []*models.Order{}
	for _, rr := range returnRequests {
		order, err := s.orderRepo.GetOrderByID(rr.OrderID)
		if err != nil {
			return nil, errors.New("unable to find order for the return request")
		}
		if order == nil {
			continue // skip if order not found
		}

		orders = append(orders, order)
	}
	return orders, nil
}

// order.userId == userCtx.id ( user )
// order.product.lenderid == userCtx.id ( lender)
