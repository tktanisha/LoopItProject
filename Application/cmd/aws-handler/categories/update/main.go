package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/services/category_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type updateCategoryRequest struct {
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
	categoryService = category_service.NewCategoryService(categoryRepo) // Adjust constructor as needed
}

func UpdateCategoryHandler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
	idStr := event.PathParameters["id"]
	categoryID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || categoryID <= 0 {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid category ID"), nil
	}

	var payload updateCategoryRequest
	if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
	}

	if payload.Price < 0 || payload.Security < 0 {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Price and security must be non-negative"), nil
	}

	if err := categoryService.UpdateCategory(categoryID, payload.Name, payload.Price, payload.Security); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Failed to update category"), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]any{
		"status":  true,
		"message": "Category updated successfully",
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(UpdateCategoryHandler)))
}
