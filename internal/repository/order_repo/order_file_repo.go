package order_repo

import (
	"errors"
	"loopit/internal/enums/order_status"
	"loopit/internal/models"
	"loopit/internal/repository/product_repo"
	"loopit/internal/storage"
)

type OrderFileRepo struct {
	filePath    string
	orders      []models.Order
	productRepo product_repo.ProductRepo
}

func NewOrderFileRepo(filePath string, productRepo product_repo.ProductRepo) (*OrderFileRepo, error) {
	orders, err := storage.ReadJSONFile[models.Order](filePath)
	if err != nil {
		return nil, err
	}
	return &OrderFileRepo{
		filePath:    filePath,
		orders:      orders,
		productRepo: productRepo,
	}, nil
}

func (r *OrderFileRepo) CreateOrder(order models.Order) error {
	order.ID = len(r.orders) + 1
	r.orders = append(r.orders, order)
	return nil
}

func (r *OrderFileRepo) UpdateOrderStatus(orderID int, newStatus string) error {
	for i, ord := range r.orders {
		if ord.ID == orderID {
			newStatusEnum, err := order_status.ParseStatus(newStatus)
			if err != nil {
				return err
			}
			r.orders[i].Status = newStatusEnum
			return nil
		}
	}
	return errors.New("order not found")
}

func (r *OrderFileRepo) GetOrderHistory(userID int, filterStatuses []string) ([]*models.Order, error) {
	filtered := []*models.Order{}

	statusMap := make(map[string]bool)
	for _, status := range filterStatuses {
		statusMap[status] = true
	}

	for _, ord := range r.orders {
		if ord.UserID != userID {
			continue
		}
		if len(filterStatuses) == 0 || statusMap[ord.Status.String()] {
			filtered = append(filtered, &ord)
		}
	}

	return filtered, nil
}

func (r *OrderFileRepo) GetLenderOrders(userID int) ([]*models.Order, error) {
	filtered := []*models.Order{}

	for _, ord := range r.orders {
		product, err := r.productRepo.FindByID(ord.ProductID)
		if err != nil {
			return nil, err
		}
		if product.Product.LenderID == userID {
			filtered = append(filtered, &ord)
		}
	}

	return filtered, nil
}

func (r *OrderFileRepo) GetOrderByID(orderID int) (*models.Order, error) {
	for _, ord := range r.orders {
		if ord.ID == orderID {
			return &ord, nil
		}
	}
	return nil, errors.New("order not found")
}

func (r *OrderFileRepo) Save() error {
	return storage.WriteJSONFile(r.filePath, r.orders)
}
