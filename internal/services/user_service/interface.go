package user_service

import "loopit/internal/models"

type UserServiceInterface interface {
	BecomeLender(user *models.UserContext) error
}
