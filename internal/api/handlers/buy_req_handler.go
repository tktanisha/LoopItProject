package handlers

import (
	"encoding/json"
	"fmt"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"loopit/internal/services/buyer_request_service"
	"loopit/internal/services/product_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	"strings"

	"net/http"
	"strconv"
)

type BuyerRequestHandler struct {
	buyerRequestService buyer_request_service.BuyerRequestServiceInterface
	productService      product_service.ProductServiceInterface
	log                 logger.LoggerInterface
}

func NewBuyerRequestHandler(buyerRequestService buyer_request_service.BuyerRequestServiceInterface, product_service product_service.ProductServiceInterface, log logger.LoggerInterface) *BuyerRequestHandler {
	return &BuyerRequestHandler{buyerRequestService: buyerRequestService, productService: product_service, log: log}
}

func (h *BuyerRequestHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("POST /buyer-requests", h.CreateBuyerRequest)
	r.HandleFunc("GET /buyer-requests", h.GetAllBuyerRequests)
	r.HandleFunc("PATCH /buyer-requests/{requestId}/update", h.UpdateBuyerRequestStatus)
}

// POST /buyer-requests
func (h *BuyerRequestHandler) CreateBuyerRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in the handler")
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}
	fmt.Println("userCtx=", userCtx)

	var payload struct {
		ProductID int `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	fmt.Println("payload=", payload)

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

// GET /buyer-requests?product_id=123&status=pending,approved
func (h *BuyerRequestHandler) GetAllBuyerRequests(w http.ResponseWriter, r *http.Request) {
	var productID *int
	if productIDStr := r.URL.Query().Get("product_id"); productIDStr != "" {
		id, err := strconv.Atoi(productIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid product_id query param", "")
			return
		}
		productID = &id
	}

	var statusFilter []string
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		statusFilter = strings.Split(statusStr, ",")
	}

	requests, err := h.buyerRequestService.GetAllBuyerRequests(productID, statusFilter)
	if err != nil {
		h.log.Error("Failed to get buyer requests: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch buyer requests", err.Error())
		return
	}

	var requestResponses []*models.BuyingRequestDto

	for _, req := range requests {
		fmt.Println("request from request array=", req)
		fmt.Println(h.productService)
		product, err := h.productService.GetProductByID(req.ProductID)
		if err != nil {
			h.log.Warning("Failed to fetch product for buyer request: " + err.Error())

			continue
		}

		requestResponses = append(requestResponses, &models.BuyingRequestDto{BuyRequest: req, Product: *product})

		fmt.Println("all request=", requestResponses)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":   true,
		"requests": requestResponses,
	})
}
