package handlers

import (
	"loopit/internal/api/router"
	"loopit/internal/services/product_service"
	"loopit/pkg/logger"
	"net/http"
)

type ProductHandler struct {
	productService product_service.ProductServiceInterface
	log            *logger.Logger
}

func NewProductHandler(productService product_service.ProductServiceInterface, log *logger.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		log:            log,
	}
}

func (h *ProductHandler) RegisterRoutes(r router.Router) {
	r.Handle("/products", http.HandlerFunc(h.GetAllProducts))

}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Fetching all products")
}
