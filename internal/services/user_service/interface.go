package user_service

import "loopit/internal/models"

type UserServiceInterface interface {
	BecomeLender(user *models.UserContext) error
	GetAllUsers(filters models.UserFilter) ([]*models.User, error)
	GetUserByID(id int) (*models.User, error)
	DeleteUserByID(id int) error
}
