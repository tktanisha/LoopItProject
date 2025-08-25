package handlers

import (
	"loopit/internal/api/router"
	"loopit/internal/services/return_request_service"
	"loopit/pkg/logger"
	"net/http"
)

type ReturnRequestHandler struct {
	returnRequestService return_request_service.ReturnRequestServiceInterface
	log                  *logger.Logger
}

func NewReturnRequestHandler(returnRequestService return_request_service.ReturnRequestServiceInterface, log *logger.Logger) *ReturnRequestHandler {
	return &ReturnRequestHandler{returnRequestService: returnRequestService, log: log}
}

func (h *ReturnRequestHandler) RegisterRoutes(r router.Router) {
	r.Handle("/return-requests", http.HandlerFunc(h.GetAllReturnRequests))
}

func (h *ReturnRequestHandler) GetAllReturnRequests(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Fetching all return requests")
}
