package user_service

import (
	"errors"
	"fmt"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
	"loopit/pkg/logger"
)

type UserService struct {
	userRepo user_repo.UserRepo
	log      *logger.Logger
}

func NewUserService(repo user_repo.UserRepo, log *logger.Logger) UserServiceInterface {
	return &UserService{userRepo: repo, log: log}
}

func (s *UserService) BecomeLender(user *models.UserContext) error {
	if user.Role == enums.RoleLender {
		s.log.Info("User is already a lender")
		return errors.New("user is already a lender")
	}

	err := s.userRepo.BecomeLender(user.ID)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to update user role to lender: %v", err))
		return err
	}

	s.log.Info("User successfully became a lender")
	return nil
}
