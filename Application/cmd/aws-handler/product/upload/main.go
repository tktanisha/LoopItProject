package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"loopit/internal/api/middleware"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var bucketName = "loopit-product-image"

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileName := event.PathParameters["fileName"]
	log.Print(fileName)
	if fileName == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "fileName is required"}, nil
	}

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-south-1"))
	if err != nil {
		log.Print("err=", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	client := s3.NewFromConfig(cfg)

	// Create presign client
	presignClient := s3.NewPresignClient(client)

	// Generate pre-signed PUT URL with ACL public-read
	presignParams := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		ACL:    types.ObjectCannedACLPublicReadWrite,
	}

	presignResult, err := presignClient.PresignPutObject(ctx, presignParams, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	response := map[string]string{
		"uploadUrl": presignResult.URL,
		"fileUrl":   fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName),
	}

	body, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(middleware.WithCORS(Handler))
}
