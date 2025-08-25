package handlers

import (
	"loopit/internal/api/router"
	"loopit/internal/services/order_service"
	"loopit/pkg/logger"
	"net/http"
)

type OrderHandler struct {
	orderService order_service.OrderServiceInterface
	log          *logger.Logger
}

func NewOrderHandler(orderService order_service.OrderServiceInterface, log *logger.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		log:          log,
	}
}

func (h *OrderHandler) RegisterRoutes(r router.Router) {
	r.Handle("/orders", http.HandlerFunc(h.GetAllOrders))

}
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Fetching all orders")

}
