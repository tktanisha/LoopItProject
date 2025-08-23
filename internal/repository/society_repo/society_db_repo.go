package society_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/models"
	"loopit/pkg/logger"
)

type SocietyDBRepo struct {
	db  *sql.DB
	log *logger.Logger
}

func NewSocietyDBRepo(db *sql.DB, log *logger.Logger) *SocietyDBRepo {
	return &SocietyDBRepo{db: db, log: log}
}

// Create inserts a new society into the database
func (r *SocietyDBRepo) Create(society models.Society) error {
	query := `
	INSERT INTO societies (name, location, pincode)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	err := r.db.QueryRow(query, society.Name, society.Location, society.Pincode).Scan(&society.ID)
	if err != nil {
		r.log.Error(fmt.Sprintf("Repo: Failed to insert society (name=%s): %v", society.Name, err))
		return err
	}
	r.log.Info(fmt.Sprintf("Repo: Society inserted successfully (id=%d, name=%s)", society.ID, society.Name))
	return nil
}

// FindAll returns all societies
func (r *SocietyDBRepo) FindAll() ([]models.Society, error) {
	rows, err := r.db.Query("SELECT id, name, location, pincode FROM societies")
	if err != nil {
		r.log.Error(fmt.Sprintf("Repo: Failed to query societies: %v", err))
		return nil, err
	}
	defer rows.Close()

	var societies []models.Society
	for rows.Next() {
		var s models.Society
		if err := rows.Scan(&s.ID, &s.Name, &s.Location, &s.Pincode); err != nil {
			r.log.Warning(fmt.Sprintf("Repo: Failed to scan society row: %v", err))
			continue
		}
		societies = append(societies, s)
	}
	r.log.Info(fmt.Sprintf("Repo: Retrieved %d societies", len(societies)))
	return societies, nil
}

// FindByID returns a society by its ID
func (r *SocietyDBRepo) FindByID(id int) (models.Society, error) {
	row := r.db.QueryRow("SELECT id, name, location, pincode FROM societies WHERE id=$1", id)
	var s models.Society
	if err := row.Scan(&s.ID, &s.Name, &s.Location, &s.Pincode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Warning(fmt.Sprintf("Repo: Society not found with id=%d", id))
			return models.Society{}, errors.New("society not found")
		}
		r.log.Error(fmt.Sprintf("Repo: Failed to scan society by id=%d: %v", id, err))
		return models.Society{}, err
	}
	r.log.Info(fmt.Sprintf("Repo: Society found (id=%d, name=%s)", s.ID, s.Name))
	return s, nil
}

// Save is a no-op for Postgres because changes are applied immediately
func (r *SocietyDBRepo) Save() error {
	return nil
}
