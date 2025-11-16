//go:generate mockgen -source=interface.go -destination=internal/repository/user_repo/mock_user_repo.go
package user_repo

import "loopit/internal/models"

type UserRepo interface {
	 FindAll(filters models.UserFilter) ([]*models.User, error)
	FindByID(userID int64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	 Create(user *models.User) error
	//BecomeLender(userID int) error
	DeleteByID(userID int64) error
	
}
