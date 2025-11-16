//go:generate mockgen -source=interface.go -destination=../../mock/mock_order_repo.go -package=mock
package order_repo

import "loopit/internal/models"

type OrderRepo interface {
	CreateOrder(order models.Order) error
	UpdateOrderStatus(orderID int64, newStatus string) error
	GetOrderHistory(userID int64, filterStatuses []string) ([]*models.Order, error)
	GetLenderOrders(userID int64) ([]*models.Order, error)
	GetOrderByID(orderID int64) (*models.Order, error)

	Save() error
}
