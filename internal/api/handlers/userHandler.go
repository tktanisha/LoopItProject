package handlers

import (
	"fmt"
	"loopit/internal/api/router"
	"loopit/internal/services/user_service"
	"loopit/pkg/logger"
	"net/http"
)

type UserHandler struct {
	userService user_service.UserServiceInterface
	log         *logger.Logger
}

func NewUserHandler(userService user_service.UserServiceInterface, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

// all handler will register their routes
func (h *UserHandler) RegisterRoutes(r router.Router) {
	r.Handle("/users", http.HandlerFunc(h.GetAllUsers)) //
	r.Handle("/users/create", http.HandlerFunc(h.CreateUser))
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Get all users")
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create user")
}
