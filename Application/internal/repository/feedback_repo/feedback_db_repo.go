package feedback_repo

import (
	"context"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FeedBackDBRepo struct {
	db *db.DynamoClient
}

func NewFeedBackDBRepo(db *db.DynamoClient) *FeedBackDBRepo {
	return &FeedBackDBRepo{db: db}
}

func (r *FeedBackDBRepo) CreateFeedback(feedback models.Feedback) error {
	feedback.ID = time.Now().UnixNano()
	feedback.CreatedAt = time.Now()

	item := map[string]types.AttributeValue{
		"pk":        &types.AttributeValueMemberS{Value: "FEEDBACK"},
		"sk":        &types.AttributeValueMemberS{Value: fmt.Sprintf("FEEDBACK#%d", feedback.ID)},
		"ID":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", feedback.ID)},
		"GivenBy":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", feedback.GivenBy)},
		"GivenTo":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", feedback.GivenTo)},
		"Text":      &types.AttributeValueMemberS{Value: feedback.Text},
		"Rating":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", feedback.Rating)},
		"CreatedAt": &types.AttributeValueMemberS{Value: feedback.CreatedAt.Format(time.RFC3339)},
	}

	_, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.db.Table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}
	return nil
}

func (r *FeedBackDBRepo) GetAllFeedbacks() ([]models.Feedback, error) {
	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "FEEDBACK"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query feedbacks: %w", err)
	}

	var feedbacks []models.Feedback
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &feedbacks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feedbacks: %w", err)
	}
	return feedbacks, nil
}

func (r *FeedBackDBRepo) Save() error {
	return nil
}
