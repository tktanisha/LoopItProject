package handlers

import (
	"loopit/internal/api/router"
	"loopit/internal/services/buyer_request_service"
	"loopit/pkg/logger"
	"net/http"
)

type BuyRequestHandler struct {
	buyRequestService buyer_request_service.BuyerRequestServiceInterface
	log               *logger.Logger
}

func NewBuyRequestHandler(buyRequestService buyer_request_service.BuyerRequestServiceInterface, log *logger.Logger) *BuyRequestHandler {
	return &BuyRequestHandler{buyRequestService: buyRequestService, log: log}
}

func (h *BuyRequestHandler) RegisterRoutes(r router.Router) {
	r.Handle("/buy-requests", http.HandlerFunc(h.GetAllBuyRequests))
}

func (h *BuyRequestHandler) GetAllBuyRequests(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Fetching all buy requests")
}
