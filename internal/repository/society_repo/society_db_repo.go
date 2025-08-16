package society_repo

import (
	"database/sql"
	"errors"
	"loopit/internal/models"
)

type SocietyDBRepo struct {
	db *sql.DB
}

func NewSocietyDBRepo(db *sql.DB) *SocietyDBRepo {
	return &SocietyDBRepo{db: db}
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
		return err
	}
	return nil
}

// FindAll returns all societies
func (r *SocietyDBRepo) FindAll() ([]models.Society, error) {
	rows, err := r.db.Query("SELECT id, name, location, pincode FROM societies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var societies []models.Society
	for rows.Next() {
		var s models.Society
		if err := rows.Scan(&s.ID, &s.Name, &s.Location, &s.Pincode); err != nil {
			continue
		}
		societies = append(societies, s)
	}

	return societies, nil
}

// FindByID returns a society by its ID
func (r *SocietyDBRepo) FindByID(id int) (models.Society, error) {
	row := r.db.QueryRow("SELECT id, name, location, pincode FROM societies WHERE id=$1", id)
	var s models.Society
	if err := row.Scan(&s.ID, &s.Name, &s.Location, &s.Pincode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Society{}, errors.New("society not found")
		}
		return models.Society{}, err
	}
	return s, nil
}

// Save is a no-op for Postgres because changes are applied immediately
func (r *SocietyDBRepo) Save() error {
	return nil
}
