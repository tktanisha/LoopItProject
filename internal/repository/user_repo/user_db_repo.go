package user_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/lender_repo"
	"loopit/pkg/logger"
	"time"
)

type UserDBRepo struct {
	db         *sql.DB
	lenderRepo lender_repo.LenderRepo
	log        *logger.Logger
}

func NewUserDBRepo(db *sql.DB, lenderRepo lender_repo.LenderRepo, log *logger.Logger) *UserDBRepo {
	return &UserDBRepo{db: db, lenderRepo: lenderRepo, log: log}
}

// FindAll returns all users
func (r *UserDBRepo) FindAll() []models.User {
	rows, err := r.db.Query("SELECT id, full_name, email, phone_number, address, password_hash, society_id, role, created_at FROM users")
	if err != nil {
		r.log.Error(fmt.Sprintf("DB query failed in FindAll: %v", err))
		return []models.User{}
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var roleStr string
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Address,
			&u.PasswordHash, &u.SocietyID, &roleStr, &u.CreatedAt); err != nil {
			r.log.Warning(fmt.Sprintf("Row scan failed in FindAll: %v", err))
			continue
		}
		u.Role, err = enums.ParseRole(roleStr)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Invalid role parsing in FindAll for user: %s, error: %v", u.Email, err))
			continue
		}
		users = append(users, u)
	}
	return users
}

// FindByID returns a user by ID
func (r *UserDBRepo) FindByID(userID int) (*models.User, error) {
	row := r.db.QueryRow("SELECT id, full_name, email, phone_number, address, password_hash, society_id, role, created_at FROM users WHERE id=$1", userID)
	var u models.User
	var roleStr string
	if err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Address,
		&u.PasswordHash, &u.SocietyID, &roleStr, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		r.log.Error(fmt.Sprintf("DB scan failed in FindByID for userID=%d, error: %v", userID, err))
		return nil, err
	}
	role, err := enums.ParseRole(roleStr)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Invalid role parsing in FindByID for userID=%d, error: %v", userID, err))
		return nil, err
	}
	u.Role = role
	return &u, nil
}

// FindByEmail returns a user by email
func (r *UserDBRepo) FindByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow("SELECT id, full_name, email, phone_number, address, password_hash, society_id, role, created_at FROM users WHERE email=$1", email)
	var u models.User
	var roleStr string
	if err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Address,
		&u.PasswordHash, &u.SocietyID, &roleStr, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		r.log.Error(fmt.Sprintf("DB scan failed in FindByEmail for email=%s, error: %v", email, err))
		return nil, err
	}
	role, err := enums.ParseRole(roleStr)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Invalid role parsing in FindByEmail for email=%s, error: %v", email, err))
		return nil, err
	}
	u.Role = role
	return &u, nil
}

// Create inserts a new user into the database
func (r *UserDBRepo) Create(user *models.User) {
	query := `
	INSERT INTO users (full_name, email, phone_number, address, password_hash, society_id, role, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`
	err := r.db.QueryRow(query, user.FullName, user.Email, user.PhoneNumber, user.Address,
		user.PasswordHash, user.SocietyID, user.Role.String(), time.Now()).Scan(&user.ID)
	if err != nil {
		r.log.Error(fmt.Sprintf("DB insert failed in Create for email=%s, error: %v", user.Email, err))
		panic(err) // TODO: replace panic with proper error handling in prod
	}
	r.log.Info(fmt.Sprintf("User inserted into DB successfully, email=%s, id=%d", user.Email, user.ID))
}

// BecomeLender updates a user's role to "lender" and creates a lender entry
func (r *UserDBRepo) BecomeLender(userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.log.Error(fmt.Sprintf("Transaction begin failed in BecomeLender for userID=%d, error: %v", userID, err))
		return err
	}

	_, err = tx.Exec("UPDATE users SET role=$1 WHERE id=$2", enums.RoleLender.String(), userID)
	if err != nil {
		tx.Rollback()
		r.log.Error(fmt.Sprintf("Failed to update user role in BecomeLender for userID=%d, error: %v", userID, err))
		return err
	}

	err = r.lenderRepo.Create(&models.Lender{
		ID:            userID,
		IsVerified:    true,
		TotalEarnings: 0.0,
	})
	if err != nil {
		tx.Rollback()
		r.log.Error(fmt.Sprintf("Failed to create lender in BecomeLender for userID=%d, error: %v", userID, err))
		return err
	}

	r.log.Info(fmt.Sprintf("User promoted to lender successfully, userID=%d", userID))
	return tx.Commit()
}

// Save is not needed for Postgres as changes are applied immediately
func (r *UserDBRepo) Save() error {
	return nil
}
