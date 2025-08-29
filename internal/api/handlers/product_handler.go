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
}

// GetAllProducts returns all products
// TODO: page and limit
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		h.log.Error("Error fetching products: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch products", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
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

	json.NewEncoder(w).Encode(map[string]interface{}{
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

	if err := h.productService.CreateProduct(&product, userCtx); err != nil {
		h.log.Error("Product creation failed: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusForbidden, "failed to create product", err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "product created successfully",
		"product": product,
	})
}
