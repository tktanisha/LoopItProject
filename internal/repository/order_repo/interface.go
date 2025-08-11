package order_repo

import "loopit/internal/models"

type OrderRepo interface {
	CreateOrder(order models.Order) error
	UpdateOrderStatus(orderID int, newStatus string) error
	GetOrderHistory(userID int, filterStatuses []string) ([]*models.Order, error)
	GetLenderOrders(userID int) ([]*models.Order, error)
	GetOrderByID(orderID int) (*models.Order, error)

	Save() error
}
