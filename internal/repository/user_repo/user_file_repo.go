package user_repo

import (
	"errors"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/storage"
)

type UserFileRepo struct {
	userFile   string
	users      []models.User // data
	lenderRepo lender_repo.LenderRepo
}

func NewUserFileRepo(userFile string, lenderRepo lender_repo.LenderRepo) (*UserFileRepo, error) {
	users, err := storage.ReadJSONFile[models.User](userFile)
	if err != nil {
		return nil, err
	}
	return &UserFileRepo{
		userFile:   userFile,
		users:      users,
		lenderRepo: lenderRepo,
	}, nil
}

func (r *UserFileRepo) FindAll() []models.User {
	return r.users
}

func (r *UserFileRepo) Create(user *models.User) {
	user.ID = len(r.users) + 1
	r.users = append(r.users, *user)
}

func (r *UserFileRepo) FindByEmail(email string) (*models.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserFileRepo) FindByID(userID int) (*models.User, error) {
	for _, u := range r.users {
		if u.ID == userID {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserFileRepo) BecomeLender(userID int) error {
	found := false

	for i := range r.users {
		if r.users[i].ID == userID {
			r.users[i].Role = enums.RoleLender
			found = true
			break
		}
	}
	if !found {
		return errors.New("user not found")
	}
	return r.lenderRepo.Create(&models.Lender{
		ID:            userID,
		IsVerified:    true,
		TotalEarnings: 0.0,
	})
}

func (r *UserFileRepo) Save() error {
	return storage.WriteJSONFile(r.userFile, r.users)
}
