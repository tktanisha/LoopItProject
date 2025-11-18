package order_service

import (
	"errors"
	"log"
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

func NewOrderService(
	OrderRepo order_repo.OrderRepo,
	returnRepo return_request_repo.ReturnRequestRepo,
	productRepo product_repo.ProductRepo,
) OrderServiceInterface {
	return &OrderService{
		orderRepo:         OrderRepo,
		returnRequestRepo: returnRepo,
		productRepo:       productRepo,
	}
}

// in use(on create order), return-requested(on create return request), returned(on mark as returned)
func (s *OrderService) UpdateOrderStatus(orderID int64, newStatus order_status.Status) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	if newStatus == order_status.Returned && order.Status != order_status.ReturnRequested {
		return errors.New("order must be in return_requested status to mark as returned")
	}

	if err := s.orderRepo.UpdateOrderStatus(orderID, newStatus.String()); err != nil {
		return err
	}

	return nil
}

// GetOrderHistory returns orders for the user based on role
func (s *OrderService) GetOrderHistory(userCtx *models.UserContext, filterStatus []order_status.Status) ([]*models.Order, error) {
	filterStatusStr := []string{}
	for _, status := range filterStatus {
		filterStatusStr = append(filterStatusStr, status.String())
	}

	orders, err := s.orderRepo.GetOrderHistory(userCtx.ID, filterStatusStr)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) GetLenderOrders(userCtx *models.UserContext) ([]*models.Order, error) {
	if userCtx.Role != enums.RoleLender {
		return nil, errors.New("only lender can get orders")
	}

	orders, err := s.orderRepo.GetLenderOrders(userCtx.ID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Lender marks order as returned (status: returned)
func (s *OrderService) MarkOrderAsReturned(orderID int64, userCtx *models.UserContext) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		log.Print("err in service1=",err)
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	product, err := s.productRepo.FindByID(order.ProductID)
	if err != nil {
		log.Print("err in service2=",err)
		return errors.New("unable to find product for the order")
	}
	if product == nil {
		return errors.New("product not found")
	}
	if product.Product.LenderID != userCtx.ID {
		return errors.New("unauthorized lender")
	}

	// returnRequests, err := s.returnRequestRepo.GetAllReturnRequests([]string{return_request_status.Approved.String()})
	// if err != nil {
	// 	s.log.Error(fmt.Sprintf("Failed to fetch return requests for order %d, error: %v", orderID, err))
	// 	return errors.New("unable to find return requests for the order")
	// }

	// isExists := false
	// for _, rr := range returnRequests {
	// 	if rr.OrderID == orderID {
	// 		isExists = true
	// 		break
	// 	}
	// }
	// if !isExists {
	// 	s.log.Warning(fmt.Sprintf("No approved return request found for order %d", orderID))
	// 	return errors.New("order has not been approved for return")
	// }

	if err := s.orderRepo.UpdateOrderStatus(orderID, order_status.Returned.String()); err != nil {
		
		return err
	}
	return nil
}

// GetAllApprovedAwaitingOrders returns orders approved for return but not yet marked returned
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
			continue
		}
		if order.Status != order_status.ReturnRequested {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}
