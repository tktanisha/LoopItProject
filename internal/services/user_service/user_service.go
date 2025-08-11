package user_service

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
)

type UserService struct {
	userRepo user_repo.UserRepo
}

func NewUserService(repo user_repo.UserRepo) UserServiceInterface {
	return &UserService{userRepo: repo}
}

func (s *UserService) BecomeLender(user *models.UserContext) error {
	if user.Role == enums.RoleLender {
		return errors.New("user is already a lender")
	}

	err := s.userRepo.BecomeLender(user.ID)
	if err != nil {
		return err
	}

	return nil
}
