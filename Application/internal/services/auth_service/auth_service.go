package auth_service

import (
	"errors"
	"fmt"
	"log"
	"time"

	//"loopit/internal/enums"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/user_repo"
	"loopit/internal/utils"
	"loopit/pkg/logger"
	//"time"
)

type AuthService struct {
	userRepo user_repo.UserRepo
	log      logger.LoggerInterface
}

func NewAuthService(repo user_repo.UserRepo) AuthServiceInterface {
	return &AuthService{
		userRepo: repo,

	}
}

func (a *AuthService) Register(user *models.User) error {

	_, err := a.userRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("user already exists")
	}
	log.Print("user not found by same email")

	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}


	user.CreatedAt = time.Now()
	user.PasswordHash = hash
	user.Role = enums.RoleUser

	err = a.userRepo.Create(user) 
	if err!= nil{
		log.Print("error after create =",err)
      return err
	}
	return nil
}

func (a *AuthService) Login(email, password string) (string, *models.User, error) {

	fmt.Println("email=", email)
	fmt.Println("password=", password)

	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		log.Print("error after repo=",err)
		return "", nil, errors.New("invalid credentials")
	}
	log.Print("coming back from repo=",user)

	fmt.Print("after calling user repo for email")

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Role.String())
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
