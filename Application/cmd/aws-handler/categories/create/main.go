package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/services/category_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type createCategoryRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Security float64 `json:"security"`
}

var categoryService category_service.CategoryServiceInterface
var categoryRepo category_repo.CategoryRepo

func init() {

	dynamo, err := db.ConnectDynamo()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DynamoDB: %v", err))
	}
	log.Print("entered in the init ")

	categoryRepo = category_repo.NewCategoryDBRepo(dynamo)
	categoryService = category_service.NewCategoryService(categoryRepo)
}

func CreateCategoryHandler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
	var payload createCategoryRequest
	if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
	}

	if payload.Price < 0 || payload.Security < 0 {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Price and security must be non-negative"), nil
	}

	if err := categoryService.CreateCategory(payload.Name, payload.Price, payload.Security); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Failed to create category"), nil
	}

	return response.LambdaResponse(http.StatusCreated, map[string]any{
		"status":  true,
		"message": "Category created successfully",
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(CreateCategoryHandler)))
}
