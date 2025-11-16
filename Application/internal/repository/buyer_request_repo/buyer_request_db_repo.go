// package buyer_request_repo

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"loopit/internal/db"
// 	"loopit/internal/enums/buyer_request_status"
// 	"loopit/internal/models"
// 	"loopit/pkg/logger"
// 	"strings"
// 	"time"

// 	"github.com/lib/pq"
// )

// type BuyerRequestDBRepo struct {
// 	db  *db.DynamoClient
// 	log logger.LoggerInterface
// }

// func NewBuyerRequestDBRepo(db *db.DynamoClient) *BuyerRequestDBRepo {
// 	return &BuyerRequestDBRepo{db: db,}
// }

// // GetAllBuyerRequests returns all buyer requests, optionally filtered by productID and/or status
// func (r *BuyerRequestDBRepo) GetAllBuyerRequests(productID *int, filterStatuses []string) ([]models.BuyingRequest, error) {
// 	query := "SELECT id, product_id, requested_by, status, created_at FROM buying_requests"
// 	args := []any{}
// 	whereClauses := []string{}
// 	argCount := 1

// 	if productID != nil {
// 		whereClauses = append(whereClauses, fmt.Sprintf("product_id = $%d", argCount))
// 		args = append(args, *productID)
// 		argCount++
// 	}

// 	if len(filterStatuses) > 0 {
// 		whereClauses = append(whereClauses, fmt.Sprintf("status = ANY($%d)", argCount))
// 		args = append(args, pq.Array(filterStatuses))
// 		argCount++
// 	}

// 	if len(whereClauses) > 0 {
// 		query += " WHERE " + strings.Join(whereClauses, " AND ")
// 	}

// 	rows, err := r.db.Query(query, args...)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("DB query failed while fetching buyer requests: %v", err))
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var requests []models.BuyingRequest
// 	for rows.Next() {
// 		var rq models.BuyingRequest
// 		var statusStr string
// 		if err := rows.Scan(&rq.ID, &rq.ProductID, &rq.RequestedBy, &statusStr, &rq.CreatedAt); err != nil {
// 			r.log.Warning(fmt.Sprintf("Failed to scan row in GetAllBuyerRequests: %v", err))
// 			continue
// 		}
// 		rq.Status, err = buyer_request_status.ParseStatus(statusStr)
// 		if err != nil {
// 			r.log.Warning(fmt.Sprintf("Invalid status '%s' found in DB row: %v", statusStr, err))
// 			continue
// 		}
// 		requests = append(requests, rq)
// 	}

// 	return requests, nil
// }

// // UpdateStatusBuyerRequest updates the status of a buyer request
// func (r *BuyerRequestDBRepo) UpdateStatusBuyerRequest(id int, newStatus string) error {
// 	result, err := r.db.Exec("UPDATE buying_requests SET status=$1 WHERE id=$2", newStatus, id)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("Failed to update buyer request %d to status '%s': %v", id, newStatus, err))
// 		return err
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		r.log.Warning(fmt.Sprintf("No buyer request found for status update, id=%d", id))
// 		return errors.New("buying request not found")
// 	}

// 	return nil
// }

// // CreateBuyerRequest inserts a new buyer request into the database
// func (r *BuyerRequestDBRepo) CreateBuyerRequest(req models.BuyingRequest) error {
// 	query := `
// 	INSERT INTO buying_requests (product_id, requested_by, status, created_at)
// 	VALUES ($1, $2, $3, $4)
// 	RETURNING id
// 	`
// 	if err := r.db.QueryRow(query, req.ProductID, req.RequestedBy, req.Status.String(), time.Now()).Scan(&req.ID); err != nil {
// 		r.log.Error(fmt.Sprintf("Failed to insert buyer request (product_id=%d, user_id=%d): %v", req.ProductID, req.RequestedBy, err))
// 		return err
// 	}
// 	return nil
// }

// // GetBuyerRequestByID fetches a buyer request by ID
// func (r *BuyerRequestDBRepo) GetBuyerRequestByID(id int) (*models.BuyingRequest, error) {
// 	row := r.db.QueryRow("SELECT id, product_id, requested_by, status, created_at FROM buying_requests WHERE id=$1", id)

// 	var req models.BuyingRequest
// 	var statusStr string
// 	if err := row.Scan(&req.ID, &req.ProductID, &req.RequestedBy, &statusStr, &req.CreatedAt); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			r.log.Warning(fmt.Sprintf("Buyer request not found by ID=%d", id))
// 			return nil, errors.New("buying request not found")
// 		}
// 		r.log.Error(fmt.Sprintf("DB scan failed in GetBuyerRequestByID for id=%d: %v", id, err))
// 		return nil, err
// 	}
// 	status, err := buyer_request_status.ParseStatus(statusStr)
// 	if err != nil {
// 		r.log.Warning(fmt.Sprintf("Invalid status '%s' found for buyer request %d", statusStr, id))
// 		return nil, err
// 	}
// 	req.Status = status

// 	return &req, nil
// }

// // Save is a no-op for Postgres
// func (r *BuyerRequestDBRepo) Save() error {
// 	return nil
// }

package buyer_request_repo

import (
	"context"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BuyerRequestDBRepo struct {
	db *db.DynamoClient
}

func NewBuyerRequestDBRepo(db *db.DynamoClient) *BuyerRequestDBRepo {
	return &BuyerRequestDBRepo{
		db: db,
	}
}

func (r *BuyerRequestDBRepo) CreateBuyerRequest(req *models.BuyingRequest) error {
	req.ID = time.Now().UnixNano()
	req.CreatedAt = time.Now()

	table := r.db.Table

	// create SK variations
	skByID := fmt.Sprintf("ID#%d", req.ID)
	skByProduct := fmt.Sprintf("BUYREQ#ID#%d", req.ID)
	skByUser := fmt.Sprintf("BUYREQ#ID#%d", req.ID)
	skByStatus := fmt.Sprintf("STATUS#%s#ID#%d", strings.ToUpper(req.Status.String()), req.ID)

	item := map[string]interface{}{
		"ID":          req.ID,
		"ProductId":   req.ProductID,
		"RequestedBy": req.RequestedBy,
		"Status":      req.Status.String(),
		"CreatedAt":   req.CreatedAt.Format(time.RFC3339),
	}

	entries := []struct {
		PK string
		SK string
	}{
		{"BUYREQUEST", skByID},
		{fmt.Sprintf("PRODUCT#%d", req.ProductID), skByProduct},
		{fmt.Sprintf("USER#%d", req.RequestedBy), skByUser},
		{"BUYREQUEST", skByStatus},
	}

	for _, e := range entries {
		keys, err := attributevalue.MarshalMap(map[string]interface{}{
			"pk": e.PK,
			"sk": e.SK,
		})

		if err != nil {
			return err
		}

		body, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}

		for k, v := range body {
			keys[k] = v
		}

		_, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      keys,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *BuyerRequestDBRepo) GetBuyerRequestByID(id int64) (*models.BuyingRequest, error) {
	pk := "BUYREQUEST"
	sk := fmt.Sprintf("ID#%d", id)

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil || out.Item == nil {
		return nil, fmt.Errorf("buyer request not found")
	}

	var req models.BuyingRequest
	if err := attributevalue.UnmarshalMap(out.Item, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

// Get all for a product
func (r *BuyerRequestDBRepo) GetAllBuyerRequests(productID *int64, statuses []string) ([]models.BuyingRequest, error) {
	var pk string
	var skPrefix string

	if productID != nil {
		pk = fmt.Sprintf("PRODUCT#%d", *productID)
		skPrefix = "BUYREQ#ID#"
	} else {
		pk = "BUYREQUEST"
		skPrefix = "STATUS#"
	}

	exprValues := map[string]types.AttributeValue{
		":pk":     &types.AttributeValueMemberS{Value: pk},
		":prefix": &types.AttributeValueMemberS{Value: skPrefix},
	}

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(r.db.Table),
		KeyConditionExpression:    aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: exprValues,
	})

	if err != nil {
		return nil, err
	}

	var reqs []models.BuyingRequest
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &reqs); err != nil {
		return nil, err
	}

	// filter by status if needed
	if len(statuses) > 0 {
		filter := map[string]bool{}
		for _, s := range statuses {
			filter[strings.ToUpper(s)] = true
		}

		tmp := []models.BuyingRequest{}
		for _, r := range reqs {
			if filter[strings.ToUpper(r.Status.String())] {
				tmp = append(tmp, r)
			}
		}
		reqs = tmp
	}

	return reqs, nil
}

// Get all for user
func (r *BuyerRequestDBRepo) GetAllBuyerRequestsByUser(userID int64) ([]models.BuyingRequest, error) {
	pk := fmt.Sprintf("USER#%d", userID)

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: "BUYREQ#ID#"},
		},
	})

	if err != nil {
		return nil, err
	}

	var reqs []models.BuyingRequest
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &reqs); err != nil {
		return nil, err
	}

	return reqs, nil
}

func (r *BuyerRequestDBRepo) UpdateStatusBuyerRequest(id int64, newStatus string) error {
	// We must update ALL copies (4 items)

	statusStr := strings.ToUpper(newStatus)

	req, err := r.GetBuyerRequestByID(id)
	if err != nil {
		return err
	}

	// delete old entries
	r.DeleteBuyerRequest(id, req.ProductID, req.RequestedBy, req.Status.String())

	// update object
	req.Status, _ = buyer_request_status.ParseStatus(statusStr)

	// recreate
	return r.CreateBuyerRequest(req)
}

func (r *BuyerRequestDBRepo) DeleteBuyerRequest(id, productID, userID int64, status string) error {
	sks := []string{
		fmt.Sprintf("ID#%d", id),
		fmt.Sprintf("BUYREQ#ID#%d", id), // product
		fmt.Sprintf("BUYREQ#ID#%d", id), // user
		fmt.Sprintf("STATUS#%s#ID#%d", strings.ToUpper(status), id),
	}

	pks := []string{
		"BUYREQUEST",
		fmt.Sprintf("PRODUCT#%d", productID),
		fmt.Sprintf("USER#%d", userID),
		"BUYREQUEST",
	}

	for i := range sks {
		_, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: pks[i]},
				"sk": &types.AttributeValueMemberS{Value: sks[i]},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *BuyerRequestDBRepo) Save() error { return nil }
