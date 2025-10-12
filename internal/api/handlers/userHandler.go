package handlers

import (
	"encoding/json"
	"fmt"
	"loopit/internal/api/router"
	"loopit/internal/constants"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/services/user_service"
	"loopit/pkg/logger"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService user_service.UserServiceInterface
	log         logger.LoggerInterface
}

func NewUserHandler(userService user_service.UserServiceInterface, log logger.LoggerInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

// all handler will register their routes
func (h *UserHandler) RegisterRoutes(r router.Router) {
	r.HandleFunc("PATCH /users/become-lender", h.BecomeLender)
	r.HandleFunc("GET /users", h.GetAllUsers)
	r.HandleFunc("DELETE /users/{id}", h.DeleteUserByID)
	r.HandleFunc("GET /users/{id}", h.GetUserByID)
}

// BecomeLender controller implementation
func (h *UserHandler) BecomeLender(w http.ResponseWriter, r *http.Request) {
	userCtxVal := r.Context().Value(constants.UserCtxKey) // replace with constants.UserCtxKey if available
	if userCtxVal == nil {
		http.Error(w, "unauthorized: user context missing", http.StatusUnauthorized)
		return
	}
	userCtx, ok := userCtxVal.(*models.UserContext)
	if !ok {
		http.Error(w, "internal error: invalid user context", http.StatusInternalServerError)
		return
	}

	h.log.Info(fmt.Sprintf("User %d attempting to become a lender", userCtx.ID))

	err := h.userService.BecomeLender(userCtx)
	if err != nil {
		h.log.Error("Error promoting user to lender: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":false,"message":"become lender failed","error":"%s"}`, err.Error())
		return
	}

	userCtx.Role = enums.RoleLender

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "User promoted to lender successfully",
		"user":    userCtx,
	})

}

// Get All Users - For Admin Use Only
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.log.Error("Error fetching users: " + err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Users fetched successfully",
		"users":   users,
	})
}

// Get User By ID - For Admin Use Only
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request: invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.log.Error("Error fetching user: " + err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "User fetched successfully",
		"user":    user,
	})
}

// Delete User By ID - For Admin Use Only
func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request: invalid user ID", http.StatusBadRequest)
		return
	}
	
	err = h.userService.DeleteUserByID(userID)

	if err != nil {
		h.log.Error("Error deleting user: " + err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{

		"status":  true,
		"message": "User deleted successfully",
	})
}
