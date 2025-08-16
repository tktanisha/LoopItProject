package return_request_repo

import (
	"database/sql"
	"errors"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"time"

	"github.com/lib/pq"
)

type ReturnRequestDBRepo struct {
	db *sql.DB
}

func NewReturnRequestDBRepo(db *sql.DB) *ReturnRequestDBRepo {
	return &ReturnRequestDBRepo{db: db}
}

// CreateReturnRequest inserts a new return request into the database
func (r *ReturnRequestDBRepo) CreateReturnRequest(req models.ReturnRequest) error {
	query := `
	INSERT INTO return_requests (order_id, status, created_at)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	err := r.db.QueryRow(query, req.OrderID, req.Status, time.Now()).Scan(&req.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateReturnRequestStatus updates the status of a return request
func (r *ReturnRequestDBRepo) UpdateReturnRequestStatus(id int, newStatus string) error {
	statusEnum, err := return_request_status.ParseStatus(newStatus)
	if err != nil {
		return err
	}

	result, err := r.db.Exec("UPDATE return_requests SET status=$1 WHERE id=$2", statusEnum, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("return request not found")
	}

	return nil
}

// GetAllReturnRequests returns all return requests, optionally filtered by status
func (r *ReturnRequestDBRepo) GetAllReturnRequests(filterStatuses []string) ([]models.ReturnRequest, error) {
	baseQuery := "SELECT id, order_id, status, created_at FROM return_requests"
	args := []interface{}{}

	if len(filterStatuses) > 0 {
		baseQuery += " WHERE status = ANY($1)"
		args = append(args, pq.Array(filterStatuses)) // requires import "github.com/lib/pq"
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.ReturnRequest
	for rows.Next() {
		var rr models.ReturnRequest
		if err := rows.Scan(&rr.ID, &rr.OrderID, &rr.Status, &rr.CreatedAt); err != nil {
			continue
		}
		requests = append(requests, rr)
	}

	return requests, nil
}

// GetReturnRequestByID returns a return request by its ID
func (r *ReturnRequestDBRepo) GetReturnRequestByID(id int) (models.ReturnRequest, error) {
	row := r.db.QueryRow("SELECT id, order_id, status, created_at FROM return_requests WHERE id=$1", id)
	var rr models.ReturnRequest
	if err := row.Scan(&rr.ID, &rr.OrderID, &rr.Status, &rr.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ReturnRequest{}, errors.New("return request not found")
		}
		return models.ReturnRequest{}, err
	}
	return rr, nil
}

// Save is a no-op for Postgres as changes are applied immediately
func (r *ReturnRequestDBRepo) Save() error {
	return nil
}
