package handlers

import (
	"encoding/json"
	"loopit/internal/api/router"
	"loopit/internal/models"
	"loopit/internal/services/auth_service"
	"loopit/internal/utils"
	"loopit/pkg/logger"

	"net/http"
)

type AuthHandler struct {
	authService auth_service.AuthServiceInterface
	log         logger.LoggerInterface
}

func NewAuthHandler(authService auth_service.AuthServiceInterface, log logger.LoggerInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

// RegisterRoutes registers all authentication-related routes
func (h *AuthHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("POST /auth/login", h.Login)
	r.HandleFunc("POST /auth/register", h.Register)
	r.HandleFunc("POST /auth/logout", h.Logout)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// Login authenticates user and returns JWT token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid invalid request payload", err.Error())
		return
	}
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.log.Warning("Login failed: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "invalid credentials", err.Error())
		return
	}

	userCtx := &models.UserContext{
		ID:   user.ID,
		Name: user.FullName,
		Role: user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  userCtx,
	})
}

// Register creates a new user and returns JWT token
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if err := utils.ValidateFullName(req.FullName); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	if err := utils.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}
	if err := utils.ValidateAddress(req.Address); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	err := h.authService.Register(&models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: req.Password,
		PhoneNumber:  req.PhoneNumber,
		Address:      req.Address,
	})
	if err != nil {
		h.log.Warning("Registration failed: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, "registration failed", err.Error())
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.log.Warning("Auto login after register failed: " + err.Error())
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "invalid credentials", err.Error())
		return
	}

	userCtx := &models.UserContext{
		ID:   user.ID,
		Name: user.FullName,
		Role: user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  userCtx,
	})
}

// Logout simply acknowledges
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "logged out successfully",
	})
}
