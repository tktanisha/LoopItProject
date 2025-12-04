package return_request_repo

import (
	"context"
	"errors"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ReturnRequestDBRepo struct {
	db *db.DynamoClient
}

func NewReturnRequestDBRepo(db *db.DynamoClient) *ReturnRequestDBRepo {
	return &ReturnRequestDBRepo{db: db}
}

func (r *ReturnRequestDBRepo) CreateReturnRequest(req models.ReturnRequest) error {
	req.ID = time.Now().UnixNano()
	req.CreatedAt = time.Now()

	item := map[string]types.AttributeValue{
		"pk":          &types.AttributeValueMemberS{Value: "RETURNREQUEST"},
		"sk":          &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", req.ID)},
		"ID":          &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.ID)},
		"OrderID":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.OrderID)},
		"RequestedBy": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.RequestedBy)},
		"Status":      &types.AttributeValueMemberS{Value: req.Status.String()},
		"CreatedAt":   &types.AttributeValueMemberS{Value: req.CreatedAt.Format(time.RFC3339)},
	}

	_, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.db.Table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to create return request: %w", err)
	}
	return nil
}

func (r *ReturnRequestDBRepo) UpdateReturnRequestStatus(id int64, newStatus string) error {
	_, err := r.db.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "RETURNREQUEST"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
		},
		UpdateExpression: aws.String("SET #s = :status"),
		ExpressionAttributeNames: map[string]string{
			"#s": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: newStatus},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update return request status: %w", err)
	}
	return nil
}

func (r *ReturnRequestDBRepo) GetAllReturnRequests(filterStatuses []string) ([]models.ReturnRequest, error) {
	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "RETURNREQUEST"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query return requests: %w", err)
	}

	type requestHelper struct {
		ID          int64     `dynamodbav:"ID"`
		OrderID     int64     `dynamodbav:"OrderID"`
		RequestedBy int64     `dynamodbav:"RequestedBy"`
		StatusStr   string    `dynamodbav:"Status"`
		CreatedAt   time.Time `dynamodbav:"CreatedAt"`
	}

	var helpers []requestHelper
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &helpers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal return requests: %w", err)
	}

	var filtered []models.ReturnRequest
	statusMap := make(map[string]bool)
	for _, s := range filterStatuses {
		statusMap[s] = true
	}

	for _, h := range helpers {
		status, err := return_request_status.ParseStatus(h.StatusStr)
		if err != nil {
			continue // skip invalid status
		}

		rr := models.ReturnRequest{
			ID:          h.ID,
			OrderID:     h.OrderID,
			RequestedBy: h.RequestedBy,
			Status:      status,
			CreatedAt:   h.CreatedAt,
		}

		if len(filterStatuses) > 0 && !statusMap[rr.Status.String()] {
			continue
		}

		filtered = append(filtered, rr)
	}

	return filtered, nil
}

func (r *ReturnRequestDBRepo) GetReturnRequestByID(id int64) (models.ReturnRequest, error) {
	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "RETURNREQUEST"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
	}

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key:       key,
	})
	if err != nil {
		return models.ReturnRequest{}, fmt.Errorf("failed to get return request: %w", err)
	}
	if out.Item == nil {
		return models.ReturnRequest{}, errors.New("return request not found")
	}

	// Helper struct for unmarshalling
	type requestHelper struct {
		ID          int64     `dynamodbav:"ID"`
		OrderID     int64     `dynamodbav:"OrderID"`
		RequestedBy int64     `dynamodbav:"RequestedBy"`
		StatusStr   string    `dynamodbav:"Status"`
		CreatedAt   time.Time `dynamodbav:"CreatedAt"`
	}

	var h requestHelper
	if err := attributevalue.UnmarshalMap(out.Item, &h); err != nil {
		return models.ReturnRequest{}, fmt.Errorf("failed to unmarshal return request: %w", err)
	}

	status, err := return_request_status.ParseStatus(h.StatusStr)
	if err != nil {
		return models.ReturnRequest{}, fmt.Errorf("invalid status value: %w", err)
	}

	rr := models.ReturnRequest{
		ID:          h.ID,
		OrderID:     h.OrderID,
		RequestedBy: h.RequestedBy,
		Status:      status,
		CreatedAt:   h.CreatedAt,
	}

	return rr, nil
}

func (r *ReturnRequestDBRepo) Save() error {
	return nil
}
