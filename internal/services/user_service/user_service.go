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
	log      logger.LoggerInterface
}

func NewUserService(repo user_repo.UserRepo, log logger.LoggerInterface) UserServiceInterface {
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


	func (s *UserService) GetAllUsers(filters models.UserFilter) ([]*models.User, error) {
	s.log.Info("Fetching all users with filters")

	users, err := s.userRepo.FindAll(filters)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to fetch users: %v", err))
		return nil, err
	}

	s.log.Info(fmt.Sprintf("Fetched %d users successfully", len(users)))
	return users, nil
}



func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to retrieve user by ID %d: %v", id, err))
		return nil, err
	}

	s.log.Info(fmt.Sprintf("Successfully retrieved user by ID %d", id))
	return user, nil
}

func (s *UserService) DeleteUserByID(id int) error {
	err := s.userRepo.DeleteByID(id)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to delete user by ID %d: %v", id, err))
		return err
	}

	s.log.Info(fmt.Sprintf("Successfully deleted user by ID %d", id))
	return nil
}
