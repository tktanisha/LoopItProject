package handlers

import (
	"encoding/json"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"loopit/internal/services/buyer_request_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"

	"net/http"
	"strconv"
)

type BuyerRequestHandler struct {
	buyerRequestService buyer_request_service.BuyerRequestServiceInterface
	log                 *logger.Logger
}

func NewBuyerRequestHandler(buyerRequestService buyer_request_service.BuyerRequestServiceInterface, log *logger.Logger) *BuyerRequestHandler {
	return &BuyerRequestHandler{buyerRequestService: buyerRequestService, log: log}
}

func (h *BuyerRequestHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("POST /buyer-requests", h.CreateBuyerRequest)
	r.HandleFunc("GET /buyer-requests", h.GetAllBuyerRequests)
	r.HandleFunc("PATCH /buyer-requests/{requestId}/status", h.UpdateBuyerRequestStatus)
}

// POST /buyer-requests
func (h *BuyerRequestHandler) CreateBuyerRequest(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	var payload struct {
		ProductID int `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := h.buyerRequestService.CreateBuyerRequest(payload.ProductID, userCtx); err != nil {
		h.log.Error("Failed to create buyer request: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to create buyer request", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "buyer request created successfully",
	})
}

// GET /buyer-requests?product_id=123&status=pending
func (h *BuyerRequestHandler) GetAllBuyerRequests(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "missing product_id query param", "")
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid product_id query param", "")
		return
	}

	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		statusStr = buyer_request_status.Pending.String()
	}
	statusEnum, err := buyer_request_status.ParseStatus(statusStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid status query param", err.Error())
		return
	}

	requests, err := h.buyerRequestService.GetAllBuyerRequestsByStatus(productID, statusEnum)
	if err != nil {
		h.log.Error("Failed to get buyer requests: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch buyer requests", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   true,
		"requests": requests,
	})
}

// PATCH /buyer-requests/{requestId}/status
func (h *BuyerRequestHandler) UpdateBuyerRequestStatus(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	reqIDStr := r.PathValue("requestId")
	reqID, err := strconv.Atoi(reqIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid buyer request ID", "")
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	statusEnum, err := buyer_request_status.ParseStatus(payload.Status)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid status value", err.Error())
		return
	}

	if err := h.buyerRequestService.UpdateBuyerRequestStatus(reqID, statusEnum, userCtx); err != nil {
		h.log.Error("Failed to update buyer request status: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to update status", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "buyer request status updated successfully",
	})
}
