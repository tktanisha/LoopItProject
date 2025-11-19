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
//
//	func (r *BuyerRequestDBRepo) Save() error {
//		return nil
//	}
package buyer_request_repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"loopit/internal/db"
	"loopit/internal/enums/buyer_request_status"
	"loopit/internal/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BuyerRequestDBRepo struct {
    db *db.DynamoClient
}

// Constructor
func NewBuyerRequestDBRepo(db *db.DynamoClient) *BuyerRequestDBRepo {
    return &BuyerRequestDBRepo{db: db}
}

// ✅ CreateBuyerRequest
func (r *BuyerRequestDBRepo) CreateBuyerRequest(req models.BuyingRequest) error {
    req.ID = time.Now().UnixNano()
    req.CreatedAt = time.Now()

    // Item for general listing
    itemGeneral := map[string]types.AttributeValue{
        "pk":         &types.AttributeValueMemberS{Value: "BUYREQUEST"},
        "sk":         &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", req.ID)},
        "ID":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.ID)},
        "ProductId":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.ProductID)},
        "RequestedBy": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.RequestedBy)},
        "Status":     &types.AttributeValueMemberS{Value: req.Status.String()},
        "CreatedAt":  &types.AttributeValueMemberS{Value: req.CreatedAt.Format(time.RFC3339)},
    }

    // Item for status-based filtering
    itemStatus := map[string]types.AttributeValue{
        "pk":         &types.AttributeValueMemberS{Value: "BUYREQUEST"},
        "sk":         &types.AttributeValueMemberS{Value: fmt.Sprintf("STATUS#%s#ID#%d", req.Status.String(), req.ID)},
        "ID":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.ID)},
        "ProductId":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.ProductID)},
        "RequestedBy": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.RequestedBy)},
        "Status":     &types.AttributeValueMemberS{Value: req.Status.String()},
        "CreatedAt":  &types.AttributeValueMemberS{Value: req.CreatedAt.Format(time.RFC3339)},
    }

    // Put both items
    _, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
        TableName: aws.String(r.db.Table),
        Item:      itemGeneral,
    })
    if err != nil {
        return fmt.Errorf("failed to create buyer request general item: %w", err)
    }

    _, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
        TableName: aws.String(r.db.Table),
        Item:      itemStatus,
    })
    if err != nil {
        return fmt.Errorf("failed to create buyer request status item: %w", err)
    }

    return nil
}

func (r *BuyerRequestDBRepo) GetAllBuyerRequests(id *int64, filterStatuses []string) ([]models.BuyingRequest, error) {
    var skPrefix string
    if len(filterStatuses) > 0 {
        skPrefix = fmt.Sprintf("STATUS#%s", filterStatuses[0])
    }
    log.Print("id and filter=",id,filterStatuses)

        var queryInput dynamodb.QueryInput
        if len(filterStatuses) > 0 {
            queryInput = dynamodb.QueryInput{
                TableName:              aws.String(r.db.Table),
                KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
                ExpressionAttributeValues: map[string]types.AttributeValue{
                    ":pk":       &types.AttributeValueMemberS{Value: "BUYREQUEST"},
                    ":skPrefix": &types.AttributeValueMemberS{Value: skPrefix},
                },
            }
        } else {
            queryInput = dynamodb.QueryInput{
                TableName:              aws.String(r.db.Table),
                KeyConditionExpression: aws.String("pk = :pk"),
                ExpressionAttributeValues: map[string]types.AttributeValue{
                    ":pk": &types.AttributeValueMemberS{Value: "BUYREQUEST"},
                },
            }
        }

        out, err := r.db.Client.Query(context.TODO(), &queryInput)
    
    if err != nil {
        log.Print("afetr query",err)
        return nil, fmt.Errorf("failed to query buyer requests: %w", err)
    }

    type requestHelper struct {
        ID          int64     `dynamodbav:"ID"`
        ProductID   int64     `dynamodbav:"ProductId"`
        RequestedBy int64     `dynamodbav:"RequestedBy"`
        StatusStr   string    `dynamodbav:"Status"`
        CreatedAt   time.Time `dynamodbav:"CreatedAt"`
    }

    var helpers []requestHelper
    if err := attributevalue.UnmarshalListOfMaps(out.Items, &helpers); err != nil {
        return nil, fmt.Errorf("failed to unmarshal buyer requests: %w", err)
    }


    var filtered []models.BuyingRequest
    for _, h := range helpers {
        status, err := buyer_request_status.ParseStatus(h.StatusStr)
        if err != nil {
            log.Print("parsing",err)
            continue // or log error
        }

        req := models.BuyingRequest{
            ID:          h.ID,
            ProductID:   h.ProductID,
            RequestedBy: h.RequestedBy,
            Status:      status,
            CreatedAt:   h.CreatedAt,
        }

        
        if id != nil && req.ProductID != *id {
            continue
        }
        if len(filterStatuses) > 0 && req.Status.String() != filterStatuses[0] {
            continue
        }

        filtered = append(filtered, req)
        log.Print("filtered= 1",filtered)
    }

    return filtered, nil
}


func (r *BuyerRequestDBRepo) UpdateStatusBuyerRequest(id int64, newStatus string) error {
    if id <= 0 {
        return fmt.Errorf("invalid id: %d", id)
    }

    ctx := context.TODO()

    // 1. Fetch existing request to get old status
    key := map[string]types.AttributeValue{
        "pk": &types.AttributeValueMemberS{Value: "BUYREQUEST"},
        "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
    }

    out, err := r.db.Client.GetItem(ctx, &dynamodb.GetItemInput{
        TableName: aws.String(r.db.Table),
        Key:       key,
    })
    if err != nil {
        return fmt.Errorf("failed to fetch buyer request: %w", err)
    }
    if out.Item == nil {
        return fmt.Errorf("buyer request not found")
    }

    type helper struct {
        Status string `dynamodbav:"Status"`
    }
    var h helper
    if err := attributevalue.UnmarshalMap(out.Item, &h); err != nil {
        return fmt.Errorf("failed to unmarshal buyer request: %w", err)
    }

    oldStatus := h.Status

    // 2. Update main item
    _, err = r.db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
        TableName: aws.String(r.db.Table),
        Key:       key,
        UpdateExpression: aws.String("SET #s = :status"),
        ExpressionAttributeNames: map[string]string{"#s": "Status"},
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":status": &types.AttributeValueMemberS{Value: newStatus},
        },
    })
    if err != nil {
        return fmt.Errorf("failed to update main buyer request: %w", err)
    }

    // 3. Delete old status index
    _, err = r.db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
        TableName: aws.String(r.db.Table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "BUYREQUEST"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("STATUS#%s#ID#%d", oldStatus, id)},
        },
    })
    if err != nil {
        return fmt.Errorf("failed to delete old status index: %w", err)
    }

    // 4. Create new status index
    newItem := map[string]types.AttributeValue{
        "pk":     &types.AttributeValueMemberS{Value: "BUYREQUEST"},
        "sk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("STATUS#%s#ID#%d", newStatus, id)},
        "Status": &types.AttributeValueMemberS{Value: newStatus},
        "ID":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
    }

    _, err = r.db.Client.PutItem(ctx, &dynamodb.PutItemInput{
        TableName: aws.String(r.db.Table),
        Item:      newItem,
    })
    if err != nil {
        return fmt.Errorf("failed to create new status index: %w", err)
    }

    return nil
}

func (r *BuyerRequestDBRepo) GetBuyerRequestByID(id int64) (*models.BuyingRequest, error) {
    // Prepare DynamoDB key
    key := map[string]types.AttributeValue{
        "pk": &types.AttributeValueMemberS{Value: "BUYREQUEST"},
        "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
    }

    // Fetch item from DynamoDB
    out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
        TableName: aws.String(r.db.Table),
        Key:       key,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get buyer request: %w", err)
    }
    if out.Item == nil {
        return nil, errors.New("buyer request not found")
    }

    // Helper struct for unmarshalling DynamoDB attributes
    type requestHelper struct {
        ID          int64     `dynamodbav:"ID"`
        ProductID   int64     `dynamodbav:"ProductId"`
        RequestedBy int64     `dynamodbav:"RequestedBy"`
        StatusStr   string    `dynamodbav:"Status"`
        CreatedAt   time.Time `dynamodbav:"CreatedAt"`
    }

    var h requestHelper
    if err := attributevalue.UnmarshalMap(out.Item, &h); err != nil {
        return nil, fmt.Errorf("failed to unmarshal buyer request: %w", err)
    }

    // Convert Status string to enum
    status, err := buyer_request_status.ParseStatus(h.StatusStr)
    if err != nil {
        return nil, fmt.Errorf("invalid status value: %w", err)
    }

    // Build final model
    req := models.BuyingRequest{
        ID:          h.ID,
        ProductID:   h.ProductID,
        RequestedBy: h.RequestedBy,
        Status:      status,
        CreatedAt:   h.CreatedAt,
    }

    return &req, nil
}

func (r *BuyerRequestDBRepo) DeleteBuyerRequest(id int64) error {
    _, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
        TableName: aws.String(r.db.Table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "BUYREQUEST"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
        },
    })
    return err
}


func (r *BuyerRequestDBRepo) Save() error {
    return nil
}