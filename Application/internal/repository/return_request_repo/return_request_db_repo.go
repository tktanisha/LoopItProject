// package return_request_repo

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"loopit/internal/db"
// 	"loopit/internal/enums/return_request_status"
// 	"loopit/internal/models"
// 	"loopit/pkg/logger"
// 	"time"

// 	"github.com/lib/pq"
// )

// type ReturnRequestDBRepo struct {
// 	db  *db.DynamoClient
// 	log logger.LoggerInterface
// }

// func NewReturnRequestDBRepo(db *db.DynamoClient, log logger.LoggerInterface) *ReturnRequestDBRepo {
// 	return &ReturnRequestDBRepo{db: db, log: log}
// }

// // CreateReturnRequest inserts a new return request into the database
// func (r *ReturnRequestDBRepo) CreateReturnRequest(req models.ReturnRequest) error {
// 	query := `
//     INSERT INTO return_requests (order_id, requested_by,status, created_at)
//     VALUES ($1, $2, $3,$4)
//     RETURNING id
//     `
// 	err := r.db.QueryRow(query, req.OrderID, req.RequestedBy,req.Status.String(), time.Now()).Scan(&req.ID)
// 	if err != nil {
// 		r.log.Error("DB error creating return request: " + err.Error())
// 		return err
// 	}
// 	return nil
// }

// // UpdateReturnRequestStatus updates the status of a return request
// func (r *ReturnRequestDBRepo) UpdateReturnRequestStatus(id int, newStatus string) error {
// 	result, err := r.db.Exec("UPDATE return_requests SET status=$1 WHERE id=$2", newStatus, id)
// 	if err != nil {
// 		r.log.Error("DB error updating return request status: " + err.Error())
// 		return err
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		r.log.Warning("No return request found to update status for, id: " + fmt.Sprint(id))
// 		return errors.New("return request not found")
// 	}

// 	return nil
// }

// // GetAllReturnRequests returns all return requests, optionally filtered by status
// func (r *ReturnRequestDBRepo) GetAllReturnRequests(filterStatuses []string) ([]models.ReturnRequest, error) {
// 	baseQuery := "SELECT id, order_id, status, created_at FROM return_requests"
// 	args := []any{}

// 	if len(filterStatuses) > 0 {
// 		baseQuery += " WHERE status = ANY($1)"
// 		args = append(args, pq.Array(filterStatuses))
// 	}

// 	rows, err := r.db.Query(baseQuery, args...)
// 	if err != nil {
// 		r.log.Error("DB error fetching all return requests: " + err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var requests []models.ReturnRequest
// 	var statusStr string
// 	for rows.Next() {
// 		var rr models.ReturnRequest
// 		if err := rows.Scan(&rr.ID, &rr.OrderID, &statusStr, &rr.CreatedAt); err != nil {
// 			continue
// 		}
// 		rr.Status, err = return_request_status.ParseStatus(statusStr)
// 		if err != nil {
// 			continue
// 		}
// 		requests = append(requests, rr)
// 	}

// 	return requests, nil
// }

// // GetReturnRequestByID returns a return request by its ID
// func (r *ReturnRequestDBRepo) GetReturnRequestByID(id int) (models.ReturnRequest, error) {
// 	row := r.db.QueryRow("SELECT id, order_id, status, created_at FROM return_requests WHERE id=$1", id)
// 	var rr models.ReturnRequest
// 	var statusStr string
// 	if err := row.Scan(&rr.ID, &rr.OrderID, &statusStr, &rr.CreatedAt); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			r.log.Warning("No return request record found in DB for id: " + fmt.Sprint(id))
// 			return models.ReturnRequest{}, errors.New("return request not found")
// 		}
// 		r.log.Error("DB error fetching return request by id: " + err.Error())
// 		return models.ReturnRequest{}, err
// 	}
// 	status, err := return_request_status.ParseStatus(statusStr)
// 	if err != nil {
// 		r.log.Error("DB error parsing status on return request: " + err.Error())
// 		return models.ReturnRequest{}, err
// 	}
// 	rr.Status = status
// 	return rr, nil
// }

// // Save is a no-op for Postgres as changes are applied immediately
// func (r *ReturnRequestDBRepo) Save() error {
// 	return nil
// }

package return_request_repo

import (
	"context"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/pkg/logger"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ReturnRequestDBRepo struct {
	db  *db.DynamoClient
	log logger.LoggerInterface
}

func NewReturnRequestDBRepo(db *db.DynamoClient, log logger.LoggerInterface) *ReturnRequestDBRepo {
	return &ReturnRequestDBRepo{db: db, log: log}
}

func (r *ReturnRequestDBRepo) CreateReturnRequest(req models.ReturnRequest) error {
	req.ID = time.Now().UnixNano()
	req.CreatedAt = time.Now()

	table := r.db.Table

	// --- KEY GENERATION ---
	skID := fmt.Sprintf("ID#%d", req.ID)

	pkOrder := fmt.Sprintf("ORDER#%d", req.OrderID)
	skOrder := fmt.Sprintf("RETURNREQ#ID#%d", req.ID)

	pkBase := "RETURNREQUEST"

	// user indexing only for pending return requests
	pkUser := fmt.Sprintf("USER#%d", req.RequestedBy)
	skUser := fmt.Sprintf("RETURNREQ#ID#%d", req.ID)

	// Body
	item := map[string]interface{}{
		"ID":          req.ID,
		"OrderID":     req.OrderID,
		"RequestedBy": req.RequestedBy,
		"Status":      req.Status.String(),
		"CreatedAt":   req.CreatedAt.Format(time.RFC3339),
	}

	records := []struct {
		PK string
		SK string
	}{
		{pkBase, skID},     // Get by ID
		{pkOrder, skOrder}, // Get all by order
	}

	// Only store in USER#index if pending
	if req.Status == return_request_status.Pending {
		records = append(records, struct{ PK, SK string }{pkUser, skUser})
	}

	for _, rec := range records {
		key, err := attributevalue.MarshalMap(map[string]string{
			"pk": rec.PK,
			"sk": rec.SK,
		})
		if err != nil {
			return err
		}

		body, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}

		for k, v := range body {
			key[k] = v
		}

		_, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      key,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReturnRequestDBRepo) GetReturnRequestByID(id int64) (models.ReturnRequest, error) {
	pk := "RETURNREQUEST"
	sk := fmt.Sprintf("ID#%d", id)

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})

	if err != nil || out.Item == nil {
		return models.ReturnRequest{}, fmt.Errorf("return request not found")
	}

	var req models.ReturnRequest
	if err := attributevalue.UnmarshalMap(out.Item, &req); err != nil {
		return models.ReturnRequest{}, err
	}

	return req, nil
}

func (r *ReturnRequestDBRepo) GetAllReturnRequestsOfOrder(orderID int64) ([]models.ReturnRequest, error) {
	pk := fmt.Sprintf("ORDER#%d", orderID)

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: "RETURNREQ#ID#"},
		},
	})

	if err != nil {
		return nil, err
	}

	var list []models.ReturnRequest
	err = attributevalue.UnmarshalListOfMaps(out.Items, &list)
	return list, err
}

func (r *ReturnRequestDBRepo) GetPendingReturnRequestsForUser(userID int64) ([]models.ReturnRequest, error) {
	pk := fmt.Sprintf("USER#%d", userID)

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: "RETURNREQ#ID#"},
		},
	})

	if err != nil {
		return nil, err
	}

	var reqs []models.ReturnRequest
	err = attributevalue.UnmarshalListOfMaps(out.Items, &reqs)
	return reqs, err
}

func (r *ReturnRequestDBRepo) UpdateReturnRequestStatus(id int64, newStatus string) error {
	req, err := r.GetReturnRequestByID(id)
	if err != nil {
		return err
	}

	// Remove old USER index only if previously pending
	if req.Status == return_request_status.Pending {
		_, _ = r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%d", req.RequestedBy)},
				"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("RETURNREQ#ID#%d", id)},
			},
		})
	}

	req.Status, _ = return_request_status.ParseStatus(newStatus)

	// recreate entire record
	return r.CreateReturnRequest(req)
}

func (r *ReturnRequestDBRepo) Save() error { return nil }
