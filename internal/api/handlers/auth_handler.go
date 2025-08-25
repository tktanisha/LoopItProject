package handlers

import (
	"loopit/internal/api/router"
	"loopit/internal/services/auth_service"
	"loopit/pkg/logger"
	"net/http"
)

type AuthHandler struct {
	authService auth_service.AuthServiceInterface
	log         *logger.Logger
}

func NewAuthHandler(authService auth_service.AuthServiceInterface, log *logger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, log: log}
}

func (h *AuthHandler) RegisterRoutes(r router.Router) {
	r.Handle("/auth/login", http.HandlerFunc(h.Login))
	r.Handle("/auth/register", http.HandlerFunc(h.Register))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.log.Info("User login")
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.log.Info("User registration")
}
