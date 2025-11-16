package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoClient struct {
	Client *dynamodb.Client
	Table  string
}

func ConnectDynamo() (*DynamoClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatalf("failed to load SDK config, %v", err)
    }
    client := dynamodb.NewFromConfig(cfg)


	tableName := "loopitProject"

	return &DynamoClient{
		Client: client,
		Table:  tableName,
	}, nil
}