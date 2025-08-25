package handlers

import (
	"loopit/internal/api/router"
	"loopit/internal/services/feedback_service"
	"loopit/pkg/logger"
	"net/http"
)

type FeedbackHandler struct {
	feedbackService feedback_service.FeedbackServiceInterface
	log             *logger.Logger
}

func NewFeedbackHandler(feedbackService feedback_service.FeedbackServiceInterface, log *logger.Logger) *FeedbackHandler {
	return &FeedbackHandler{feedbackService: feedbackService, log: log}
}

func (h *FeedbackHandler) RegisterRoutes(r router.Router) {
	r.Handle("/feedback", http.HandlerFunc(h.GetAllFeedbacks))

}

func (h *FeedbackHandler) GetAllFeedbacks(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Fetching all feedbacks")
}
