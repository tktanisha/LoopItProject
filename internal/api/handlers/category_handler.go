package handlers

import (
	"encoding/json"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/models"
	"loopit/internal/services/category_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	"net/http"
)

type CategoryHandler struct {
	log     *logger.Logger
	service category_service.CategoryServiceInterface
}

func NewCategoryHandler(service category_service.CategoryServiceInterface, log *logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		log:     log,
		service: service,
	}
}

// RegisterRoutes attaches category routes to given router
func (h *CategoryHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("POST /categories", h.CreateCategory)
	r.HandleFunc("GET /categories", h.GetAllCategories)
}

// CreateCategory creates a new category; requires authenticated user
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Auth check
	userCtxVal := r.Context().Value(constants.UserCtxKey)
	if userCtxVal == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}
	_, ok := userCtxVal.(*models.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "internal error", "invalid user context")
		return
	}

	var payload struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Security float64 `json:"security"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := h.service.CreateCategory(payload.Name, payload.Price, payload.Security); err != nil {
		h.log.Error("Failed to create category: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to create category", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "category created successfully",
	})
}

// GetAllCategories returns list of all categories
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		h.log.Error("Failed to fetch categories: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch categories", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     true,
		"categories": categories,
	})
}
