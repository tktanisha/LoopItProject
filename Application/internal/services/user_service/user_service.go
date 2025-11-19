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


func (s *UserService) GetAllUsers(filters models.UserFilter) ([]*models.User, error) {

	users, err := s.userRepo.FindAll(filters)
	if err != nil {
		return nil, err
	}

	return users, nil
}



func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUserByID(id int64) error {
	err := s.userRepo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}
