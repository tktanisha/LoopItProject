package main

import (
	"context"
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

func GetAllCategoriesHandler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
	categories, err := categoryService.GetAllCategories()
	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, nil, "Failed to fetch categories"), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]any{
		"status":     true,
		"categories": categories,
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(GetAllCategoriesHandler)))
}
