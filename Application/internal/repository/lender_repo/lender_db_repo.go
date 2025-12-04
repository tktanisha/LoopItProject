package lender_repo

import (
	"context"
	"fmt"
	"loopit/internal/db"
	"loopit/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LenderDBRepo struct {
	db *db.DynamoClient
}

func NewLenderDBRepo(db *db.DynamoClient) *LenderDBRepo {
	return &LenderDBRepo{db: db}
}

func (r *LenderDBRepo) Create(lender *models.Lender) error {
	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LENDER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", lender.ID)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to check lender existence: %w", err)
	}
	if out.Item != nil {
		return nil
	}

	item := map[string]types.AttributeValue{
		"pk":            &types.AttributeValueMemberS{Value: "LENDER"},
		"sk":            &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", lender.ID)},
		"ID":            &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", lender.ID)},
		"IsVerified":    &types.AttributeValueMemberBOOL{Value: lender.IsVerified},
		"TotalEarnings": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", lender.TotalEarnings)},
	}

	_, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.db.Table),
		Item:      item,
	})
	return err
}

func (r *LenderDBRepo) Save() error {
	return nil
}
