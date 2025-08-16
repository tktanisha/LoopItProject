package auth_service

import (
	"errors"
	"fmt"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	"time"
)

type AuthService struct {
	userRepo user_repo.UserRepo
	log      *logger.Logger
}

func NewAuthService(repo user_repo.UserRepo, log *logger.Logger) AuthServiceInterface {
	return &AuthService{
		userRepo: repo,
		log:      log,
	}
}

func (a *AuthService) Register(user *models.User) error {
	a.log.Info(fmt.Sprintf("Register attempt for email: %s", user.Email))

	_, err := a.userRepo.FindByEmail(user.Email)
	if err == nil {
		a.log.Warning(fmt.Sprintf("User already exists with email: %s", user.Email))
		return errors.New("user already exists")
	}

	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		a.log.Error(fmt.Sprintf("Password hashing failed for email: %s, error: %v", user.Email, err))
		return err
	}

	user.CreatedAt = time.Now()
	user.PasswordHash = hash
	user.Role = enums.RoleUser

	a.userRepo.Create(user)
	a.log.Info(fmt.Sprintf("User registered successfully: %s", user.Email))
	return nil
}

func (a *AuthService) Login(email, password string) (string, *models.User, error) {
	a.log.Info(fmt.Sprintf("Login attempt for email: %s", email))

	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Login failed, user not found: %s", email))
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		a.log.Warning(fmt.Sprintf("Invalid password for email: %s", email))
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		a.log.Error(fmt.Sprintf("JWT generation failed for user ID: %d, error: %v", user.ID, err))
		return "", nil, err
	}

	a.log.Info(fmt.Sprintf("Login successful for email: %s", email))
	return token, user, nil
}
