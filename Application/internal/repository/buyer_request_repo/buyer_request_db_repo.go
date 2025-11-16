package buyer_request_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"loopit/pkg/logger"
	"strings"
	"time"

	"github.com/lib/pq"
)

type BuyerRequestDBRepo struct {
	db  *db.DynamoClient
	log logger.LoggerInterface
}

func NewBuyerRequestDBRepo(db *db.DynamoClient) *BuyerRequestDBRepo {
	return &BuyerRequestDBRepo{db: db,}
}

// GetAllBuyerRequests returns all buyer requests, optionally filtered by productID and/or status
func (r *BuyerRequestDBRepo) GetAllBuyerRequests(productID *int, filterStatuses []string) ([]models.BuyingRequest, error) {
	query := "SELECT id, product_id, requested_by, status, created_at FROM buying_requests"
	args := []any{}
	whereClauses := []string{}
	argCount := 1

	if productID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("product_id = $%d", argCount))
		args = append(args, *productID)
		argCount++
	}

	if len(filterStatuses) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("status = ANY($%d)", argCount))
		args = append(args, pq.Array(filterStatuses))
		argCount++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.log.Error(fmt.Sprintf("DB query failed while fetching buyer requests: %v", err))
		return nil, err
	}
	defer rows.Close()

	var requests []models.BuyingRequest
	for rows.Next() {
		var rq models.BuyingRequest
		var statusStr string
		if err := rows.Scan(&rq.ID, &rq.ProductID, &rq.RequestedBy, &statusStr, &rq.CreatedAt); err != nil {
			r.log.Warning(fmt.Sprintf("Failed to scan row in GetAllBuyerRequests: %v", err))
			continue
		}
		rq.Status, err = buyer_request_status.ParseStatus(statusStr)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Invalid status '%s' found in DB row: %v", statusStr, err))
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
		r.log.Error(fmt.Sprintf("Failed to update buyer request %d to status '%s': %v", id, newStatus, err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		r.log.Warning(fmt.Sprintf("No buyer request found for status update, id=%d", id))
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
	if err := r.db.QueryRow(query, req.ProductID, req.RequestedBy, req.Status.String(), time.Now()).Scan(&req.ID); err != nil {
		r.log.Error(fmt.Sprintf("Failed to insert buyer request (product_id=%d, user_id=%d): %v", req.ProductID, req.RequestedBy, err))
		return err
	}
	return nil
}

// GetBuyerRequestByID fetches a buyer request by ID
func (r *BuyerRequestDBRepo) GetBuyerRequestByID(id int) (*models.BuyingRequest, error) {
	row := r.db.QueryRow("SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE id=$1", id)

	var req models.BuyingRequest
	var statusStr string
	if err := row.Scan(&req.ID, &req.ProductID, &req.RequestedBy, &statusStr, &req.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Warning(fmt.Sprintf("Buyer request not found by ID=%d", id))
			return nil, errors.New("buying request not found")
		}
		r.log.Error(fmt.Sprintf("DB scan failed in GetBuyerRequestByID for id=%d: %v", id, err))
		return nil, err
	}
	status, err := buyer_request_status.ParseStatus(statusStr)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Invalid status '%s' found for buyer request %d", statusStr, id))
		return nil, err
	}
	req.Status = status

	return &req, nil
}

// Save is a no-op for Postgres
func (r *BuyerRequestDBRepo) Save() error {
	return nil
}
