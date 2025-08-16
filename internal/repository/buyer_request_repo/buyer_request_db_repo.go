package buyer_request_repo

import (
	"database/sql"
	"errors"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"time"

	"github.com/lib/pq"
)

type BuyerRequestDBRepo struct {
	db *sql.DB
}

func NewBuyerRequestDBRepo(db *sql.DB) *BuyerRequestDBRepo {
	return &BuyerRequestDBRepo{db: db}
}

// GetAllBuyerRequests returns all buyer requests, optionally filtered by status
func (r *BuyerRequestDBRepo) GetAllBuyerRequests(filterStatuses []string) ([]models.BuyingRequest, error) {
	query := "SELECT id, product_id, requested_by, status, created_at FROM buying_requests"
	args := []interface{}{}

	if len(filterStatuses) > 0 {
		query += " WHERE status = ANY($1)"
		args = append(args, pq.Array(filterStatuses))
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.BuyingRequest
	var statusStr string
	for rows.Next() {
		var rq models.BuyingRequest
		if err := rows.Scan(&rq.ID, &rq.ProductID, &rq.RequestedBy, &statusStr, &rq.CreatedAt); err != nil {
			continue
		}
		rq.Status, err = buyer_request_status.ParseStatus(statusStr)
		if err != nil {
			continue
		}
		requests = append(requests, rq)
	}

	return requests, nil
}

// UpdateStatusBuyerRequest updates the status of a buyer request
func (r *BuyerRequestDBRepo) UpdateStatusBuyerRequest(id int, newStatus string) error {
	result, err := r.db.Exec("UPDATE buying_requests SET status=$1 WHERE id=$2", newStatus, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("buying request not found")
	}

	return nil
}

// CreateBuyerRequest inserts a new buyer request into the database
func (r *BuyerRequestDBRepo) CreateBuyerRequest(req models.BuyingRequest) error {
	query := `
	INSERT INTO buying_requests (product_id, requested_by, status, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	return r.db.QueryRow(query, req.ProductID, req.RequestedBy, req.Status.String(), time.Now()).Scan(&req.ID)
}

// GetBuyerRequestByID fetches a buyer request by ID
func (r *BuyerRequestDBRepo) GetBuyerRequestByID(id int) (*models.BuyingRequest, error) {
	row := r.db.QueryRow("SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE id=$1", id)

	var req models.BuyingRequest
	var statusStr string
	if err := row.Scan(&req.ID, &req.ProductID, &req.RequestedBy, &statusStr, &req.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("buying request not found")
		}
		return nil, err
	}
	status, err := buyer_request_status.ParseStatus(statusStr)

	if err != nil {
		return nil, err
	}
	req.Status = status

	return &req, nil
}

// Save is a no-op for Postgres
func (r *BuyerRequestDBRepo) Save() error {
	return nil
}
