package handlers

import (
	"encoding/json"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/models"
	"loopit/internal/services/feedback_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"

	"net/http"
)

type FeedbackHandler struct {
	feedbackService feedback_service.FeedbackServiceInterface
	log             logger.LoggerInterface
}

func NewFeedbackHandler(feedbackService feedback_service.FeedbackServiceInterface, log logger.LoggerInterface) *FeedbackHandler {
	return &FeedbackHandler{feedbackService: feedbackService, log: log}
}

func (h *FeedbackHandler) RegisterRoutes(r router.Router) {

	r.HandleFunc("POST /feedbacks", h.GiveFeedback)
	r.HandleFunc("GET /feedbacks/given", h.GetAllGivenFeedbacks)
	r.HandleFunc("GET /feedbacks/received", h.GetAllReceivedFeedbacks)
}

// POST /feedback
func (h *FeedbackHandler) GiveFeedback(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	var payload struct {
		OrderID      int    `json:"order_id"`
		FeedbackText string `json:"feedback_text"`
		Rating       int    `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := utils.ValidateOrderID(payload.OrderID); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	if err := utils.ValidateFeedbackText(payload.FeedbackText); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	if err := utils.ValidateRating(payload.Rating); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := h.feedbackService.GiveFeedback(payload.OrderID, payload.FeedbackText, payload.Rating, userCtx); err != nil {
		h.log.Error("Failed to give feedback: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to give feedback", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Feedback given successfully",
	})
}

// GET /feedback/given
func (h *FeedbackHandler) GetAllGivenFeedbacks(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	feedbacks, err := h.feedbackService.GetAllGivenFeedbacks(userCtx)
	if err != nil {
		h.log.Error("Failed to fetch given feedbacks: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch given feedbacks", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    true,
		"feedbacks": feedbacks,
	})
}

// GET /feedback/received
func (h *FeedbackHandler) GetAllReceivedFeedbacks(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	feedbacks, err := h.feedbackService.GetAllReceivedFeedbacks(userCtx)
	if err != nil {
		h.log.Error("Failed to fetch received feedbacks: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch received feedbacks", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    true,
		"feedbacks": feedbacks,
	})
}
