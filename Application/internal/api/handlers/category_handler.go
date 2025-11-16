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
	"strconv"
)

type CategoryHandler struct {
	log     logger.LoggerInterface
	service category_service.CategoryServiceInterface
}

func NewCategoryHandler(service category_service.CategoryServiceInterface, log logger.LoggerInterface) *CategoryHandler {
	return &CategoryHandler{
		log:     log,
		service: service,
	}
}

// RegisterRoutes attaches category routes to given router
func (h *CategoryHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("POST /categories", h.CreateCategory)
	r.HandleFunc("GET /categories", h.GetAllCategories)
	r.HandleFunc("PUT /categories/{id}", h.UpdateCategory)
	r.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)
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

	if payload.Price < 0 || payload.Security < 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", "price and security must be non-negative")
		return
	}

	if err := h.service.CreateCategory(payload.Name, payload.Price, payload.Security); err != nil {
		h.log.Error("Failed to create category: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to create category", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
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
	json.NewEncoder(w).Encode(map[string]any{
		"status":     true,
		"categories": categories,
	})
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	categoryID, err := strconv.Atoi(idStr)

	if err != nil || categoryID <= 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid category ID", "category ID must be a positive integer")
		return
	}

	// use category model for payload
	var payload struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Security float64 `json:"security"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	if payload.Price < 0 || payload.Security < 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", "price and security must be non-negative")
		return
	}
	if err := h.service.UpdateCategory(categoryID, payload.Name, payload.Price, payload.Security); err != nil {
		h.log.Error("Failed to update category: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to update category", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "category updated successfully",
	})
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil || categoryID <= 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid category ID", "category ID must be a positive integer")
		return
	}
	if err := h.service.DeleteCategory(categoryID); err != nil {
		h.log.Error("Failed to delete category: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "failed to delete category", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "category deleted successfully",
	})
}
