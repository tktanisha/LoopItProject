package user_repo

import "loopit/internal/models"

type UserRepo interface {
	FindAll() []models.User
	FindByID(userID int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User)
	BecomeLender(userID int) error
	Save() error
}
