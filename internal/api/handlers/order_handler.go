package handlers

import (
	"encoding/json"
	"fmt"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/enums"
	"loopit/internal/enums/order_status"
	"loopit/internal/initializer"
	"loopit/internal/models"
	"loopit/internal/services/order_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	"net/http"
)

type OrderHandler struct {
	orderService order_service.OrderServiceInterface
	log          logger.LoggerInterface
}

func NewOrderHandler(orderService order_service.OrderServiceInterface, log logger.LoggerInterface) *OrderHandler {
	return &OrderHandler{orderService: orderService, log: log}
}

func (h *OrderHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("GET /orders/history", h.GetOrderHistory)
	r.HandleFunc("PATCH /orders/{orderId}/return", h.MarkOrderAsReturned)
	r.HandleFunc("GET /orders/approved-awaiting", h.GetAllApprovedAwaitingOrders)
	r.HandleFunc("GET /orders/lender", h.GetLenderOrders)
}

// GET /orders/history?status=APPROVED
func (h *OrderHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	statusStr := r.URL.Query().Get("status")
	var filterStatus []order_status.Status
	if statusStr != "" {
		st, err := order_status.ParseStatus(statusStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid status filter", err.Error())
			return
		}
		filterStatus = append(filterStatus, st)
	}

	orders, err := initializer.OrderService.GetOrderHistory(userCtx, filterStatus)
	if err != nil {
		h.log.Error("Failed to fetch order history: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch order history", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"orders": orders,
	})
}

// PATCH /orders/{orderId}/updateStatus
func (h *OrderHandler) MarkOrderAsReturned(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	orderIDStr := r.PathValue("orderId")
	var orderID int
	if _, err := fmt.Sscanf(orderIDStr, "%d", &orderID); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid order id", err.Error())
		return
	}

	if err := initializer.OrderService.MarkOrderAsReturned(orderID, userCtx); err != nil {
		h.log.Error("Failed to update order status: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to update order status", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "order status updated successfully",
	})
}

// GET /orders/approved-awaiting
func (h *OrderHandler) GetAllApprovedAwaitingOrders(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}
	if userCtx.Role != enums.RoleLender {
		utils.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "only lenders can view approved awaiting orders")
		return
	}

	orders, err := initializer.OrderService.GetAllApprovedAwaitingOrders(userCtx)
	if err != nil {
		h.log.Error("Failed to fetch approved awaiting orders: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch approved awaiting orders", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"orders": orders,
	})
}

// GET /orders/lender
func (h *OrderHandler) GetLenderOrders(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}
	if userCtx.Role != enums.RoleLender {
		utils.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "only lenders can view their orders")
		return
	}

	orders, err := initializer.OrderService.GetLenderOrders(userCtx)
	if err != nil {
		h.log.Error("Failed to fetch lender orders: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch lender orders", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"orders": orders,
	})
}
