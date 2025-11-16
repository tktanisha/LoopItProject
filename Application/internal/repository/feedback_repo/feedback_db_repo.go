// package feedback_repo

// import (
// 	"fmt"
// 	"loopit/internal/db"
// 	"loopit/internal/models"
// 	"time"
// )

// type FeedBackDBRepo struct {
// 	db  *db.DynamoClient

// }

// func NewFeedBackDBRepo(db *db.DynamoClient) *FeedBackDBRepo {
// 	return &FeedBackDBRepo{db: db}
// }

// // CreateFeedback inserts a new feedback into the database
// func (r *FeedBackDBRepo) CreateFeedback(feedback models.Feedback) error {
// 	query := `
//     INSERT INTO feedbacks (given_by, given_to, text, rating, created_at)
//     VALUES ($1, $2, $3, $4, $5)
//     RETURNING id
//     `
// 	err := r.db.QueryRow(query, feedback.GivenBy, feedback.GivenTo, feedback.Text, feedback.Rating, time.Now()).Scan(&feedback.ID)
// 	if err != nil && r.log != nil {
// 		r.log.Error(fmt.Sprintf("Repo: DB error creating feedback by user %d for user %d: %v", feedback.GivenBy, feedback.GivenTo, err))
// 	}
// 	return err
// }

// // GetAllFeedbacks returns all feedbacks from the database
// func (r *FeedBackDBRepo) GetAllFeedbacks() ([]models.Feedback, error) {
// 	rows, err := r.db.Query("SELECT id, given_by, given_to, text, rating, created_at FROM feedbacks")
// 	if err != nil {
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("Repo: DB error fetching all feedbacks: %v", err))
// 		}
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var feedbacks []models.Feedback
// 	for rows.Next() {
// 		var f models.Feedback
// 		if err := rows.Scan(&f.ID, &f.GivenBy, &f.GivenTo, &f.Text, &f.Rating, &f.CreatedAt); err != nil {
// 			if r.log != nil {
// 				r.log.Warning(fmt.Sprintf("Repo: Could not scan feedback row: %v", err))
// 			}
// 			continue
// 		}
// 		feedbacks = append(feedbacks, f)
// 	}

// 	if len(feedbacks) == 0 {
// 		return nil, nil
// 	}

// 	return feedbacks, nil
// }

// // Save is a no-op for Postgres
//
//	func (r *FeedBackDBRepo) Save() error {
//		return nil
//	}
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

// CreateFeedback inserts a new feedback into DynamoDB
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