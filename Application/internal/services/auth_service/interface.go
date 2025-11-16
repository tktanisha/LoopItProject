package auth_service

import "loopit/internal/models"

type AuthServiceInterface interface {
	Register(user *models.User) error
	Login(email, password string) (string, *models.User, error)
}
