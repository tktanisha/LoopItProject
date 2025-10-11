package handlers

import (
	"encoding/json"
	"fmt"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/internal/services/return_request_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"

	"net/http"
)

type ReturnRequestHandler struct {
	returnRequestService return_request_service.ReturnRequestServiceInterface
	log                  logger.LoggerInterface
}

func NewReturnRequestHandler(svc return_request_service.ReturnRequestServiceInterface, log logger.LoggerInterface) *ReturnRequestHandler {
	return &ReturnRequestHandler{returnRequestService: svc, log: log}
}

func (h *ReturnRequestHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("POST /return-requests", h.CreateReturnRequest)
	r.HandleFunc("GET /return-requests", h.GetPendingReturnRequests)
	r.HandleFunc("PATCH /return-requests/{requestId}/update", h.UpdateReturnRequestStatus)
}

func (h *ReturnRequestHandler) GetPendingReturnRequests(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	requests, err := h.returnRequestService.GetPendingReturnRequests(userCtx.ID)
	fmt.Println("handler=",requests)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "could not fetch return requests", err.Error())
		return
	}
    
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   true,
		"requests": requests,
	})
}

func (h *ReturnRequestHandler) CreateReturnRequest(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	var payload struct {
		OrderID int `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := h.returnRequestService.CreateReturnRequest(userCtx.ID, payload.OrderID); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "could not create return request", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Return request created successfully",
	})
}

func (h *ReturnRequestHandler) UpdateReturnRequestStatus(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	reqIDStr := r.PathValue("requestId")
	var reqID int
	if _, err := fmt.Sscanf(reqIDStr, "%d", &reqID); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request id", err.Error())
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	// convert string → enum
	newStatus, err := return_request_status.ParseStatus(payload.Status)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid status value", err.Error())
		return
	}

	if err := h.returnRequestService.UpdateReturnRequestStatus(userCtx.ID, reqID, newStatus); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "could not update return request", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Return request status updated successfully",
	})
}
