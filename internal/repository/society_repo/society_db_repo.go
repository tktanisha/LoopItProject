package society_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/pkg/logger"
)

type SocietyDBRepo struct {
	db  db.DatabaseInterface
	log logger.LoggerInterface
}

func NewSocietyDBRepo(db db.DatabaseInterface, log logger.LoggerInterface) *SocietyDBRepo {
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

// update modifies an existing society in the database
func (r *SocietyDBRepo) Update(society models.Society) error {
	query := `UPDATE societies SET name=$1, location=$2, pincode=$3 WHERE id=$4`
	result, err := r.db.Exec(query, society.Name, society.Location, society.Pincode, society.ID)
	if err != nil {
		r.log.Error(fmt.Sprintf("Repo: Failed to update society (id=%d): %v", society.ID, err))
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error(fmt.Sprintf("Repo: Failed to get rows affected for society update (id=%d): %v", society.ID, err))
		return err
	}
	if rowsAffected == 0 {
		r.log.Warning(fmt.Sprintf("Repo: No society found to update with id=%d", society.ID))
		return errors.New("society not found")
	}
	r.log.Info(fmt.Sprintf("Repo: Society updated successfully (id=%d)", society.ID))
	return nil
}

// Delete removes a society from the database by its ID
func (r *SocietyDBRepo) Delete(id int) error {
	query := `DELETE FROM societies WHERE id=$1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		r.log.Error(fmt.Sprintf("Repo: Failed to delete society (id=%d): %v", id, err))
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error(fmt.Sprintf("Repo: Failed to get rows affected for society delete (id=%d): %v", id, err))
		return err
	}
	if rowsAffected == 0 {
		r.log.Warning(fmt.Sprintf("Repo: No society found to delete with id=%d", id))
		return errors.New("society not found")
	}
	r.log.Info(fmt.Sprintf("Repo: Society deleted successfully (id=%d)", id))
	return nil
}

// Save is a no-op for Postgres because changes are applied immediately
func (r *SocietyDBRepo) Save() error {
	return nil
}
