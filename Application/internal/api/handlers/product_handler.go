package handlers

import (
	"encoding/json"
	"fmt"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/models"
	"loopit/internal/services/product_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"

	"net/http"
	"strconv"
)

type ProductHandler struct {
	productService product_service.ProductServiceInterface
	log            logger.LoggerInterface
}

func NewProductHandler(svc product_service.ProductServiceInterface, log logger.LoggerInterface) *ProductHandler {
	return &ProductHandler{productService: svc, log: log}
}

func (h *ProductHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("GET /products", h.GetAllProducts)
	r.HandleFunc("GET /products/{id}", h.GetProductByID)
	r.HandleFunc("POST /products/create", h.CreateProduct) // protected route (lender only)
	r.HandleFunc("PUT /products/{id}/update", h.UpdateProduct) // protected route (lender only)
	r.HandleFunc("DELETE /products/{id}/delete", h.DeleteProduct) // protected route (lender only)
}

// GetAllProducts returns all products
// TODO: page and limit
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	query := r.URL.Query()
	search := query.Get("search")
	lenderID := query.Get("lender_id")
	categoryID := query.Get("category_id")
	isAvailable := query.Get("is_available")

	filters := models.ProductFilter{
		Search:      search,
		LenderID:    lenderID,
		CategoryID:  categoryID,
		IsAvailable: isAvailable,
	}

	products, err := h.productService.GetAllProducts(filters)
	if err != nil {
		h.log.Error("Error fetching products: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch products", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":   true,
		"products": products,
	})
}


// GetProductByID fetches product details by ID (from path param)
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid product id", fmt.Sprintf("id=%s", idStr))
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		h.log.Error("Product not found: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusNotFound, "product not found", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"product": product,
	})
}

// CreateProduct creates a new product, requires user context (lender)
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	// Validation before processing
	if err := utils.ValidateProduct(&product); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := h.productService.CreateProduct(&product, userCtx); err != nil {
		h.log.Error("Product creation failed: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusForbidden, "failed to create product", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "product created successfully",
		"product": product,
	})
}

// UpdateProduct updates an existing product, requires user context (lender)
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid product id", fmt.Sprintf("id=%s", idStr))
		return
	}
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	// Validation before processing
	if err := utils.ValidateProduct(&product); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := h.productService.UpdateProduct(id, product.Name, product.Description, product.CategoryID, userCtx); err != nil {
		h.log.Error("Product update failed: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusForbidden, "failed to update product", err.Error())
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "product updated successfully",
		"product": product,
	})
}

// DeleteProduct deletes a product by ID, requires user context (lender)
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value(constants.UserCtxKey).(*models.UserContext)
	if !ok || userCtx == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "user context missing")
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid product id", fmt.Sprintf("id=%s", idStr))
		return
	}
	if err := h.productService.DeleteProduct(id, userCtx); err != nil {
		h.log.Error("Product deletion failed: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusForbidden, "failed to delete product", err.Error())
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "product deleted successfully",
	})
}