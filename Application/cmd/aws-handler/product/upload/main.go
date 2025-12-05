// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"loopit/internal/api/middleware"
// 	"time"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/aws/aws-sdk-go-v2/service/s3/types"
// )

// var bucketName = "loopit-product-image"

// func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	fileName := event.PathParameters["fileName"]
// 	log.Print(fileName)
// 	if fileName == "" {
// 		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "fileName is required"}, nil
// 	}

// 	// Load AWS config
// 	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-south-1"))
// 	if err != nil {
// 		log.Print("err=", err)
// 		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
// 	}

// 	client := s3.NewFromConfig(cfg)

// 	// Create presign client
// 	presignClient := s3.NewPresignClient(client)

// 	// Generate pre-signed PUT URL with ACL public-read
// 	presignParams := &s3.PutObjectInput{
// 		Bucket:      &bucketName,
// 		Key:         &fileName,
// 		ACL:         types.ObjectCannedACLPublicRead,
// 		ContentType: aws.String("image/jpeg"),
// 	}
// 	log.Print("params=", presignParams)

// 	presignResult, err := presignClient.PresignPutObject(ctx, presignParams, s3.WithPresignExpires(15*time.Minute))
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil

// 	}

// 	log.Print("result=", presignResult)

// 	response := map[string]string{
// 		"uploadUrl": presignResult.URL,
// 		"fileUrl":   fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName),
// 	}
// 	log.Print("response=", response)

// 	body, _ := json.Marshal(response)
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: 200,
// 		Headers:    map[string]string{"Content-Type": "application/json"},
// 		Body:       string(body),
// 	}, nil
// }

// func main() {
// 	lambda.Start(middleware.WithCORS(Handler))
// }

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"loopit/internal/api/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var bucketName = "loopit-product-image"

type UploadRequest struct {
	FileName    string `json:"fileName"`
	FileType    string `json:"fileType"`
	FileContent string `json:"fileContent"`
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var req UploadRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
	}

	if req.FileName == "" || req.FileType == "" || req.FileContent == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "fileName, fileType, and fileContent are required",
		}, nil
	}

	fileBytes, err := base64.StdEncoding.DecodeString(req.FileContent)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid Base64 fileContent"}, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-south-1"))
	if err != nil {
		log.Println("Config error:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	s3Client := s3.NewFromConfig(cfg)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(req.FileName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(req.FileType),
		// ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		log.Println("S3 upload error:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	fileUrl := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, req.FileName)

	response := map[string]string{
		"fileUrl": fileUrl,
	}

	respJson, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(respJson),
	}, nil
}

func main() {
	lambda.Start(middleware.WithCORS(Handler))
}
