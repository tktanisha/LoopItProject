package auth_service

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
	"loopit/internal/utils"
	"time"
)

type AuthService struct {
	userRepo user_repo.UserRepo
}

func NewAuthService(repo user_repo.UserRepo) AuthServiceInterface {
	return &AuthService{userRepo: repo}
}

func (a *AuthService) Register(user *models.User) error {
	_, err := a.userRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("user already exists")
	}

	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}

	user.CreatedAt = time.Now()
	user.PasswordHash = hash
	user.Role = enums.RoleUser // Default role for new users
	a.userRepo.Create(user)
	return nil
}

func (a *AuthService) Login(email, password string) (string, *models.User, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
