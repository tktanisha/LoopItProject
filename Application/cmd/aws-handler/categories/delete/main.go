package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"loopit/internal/db"
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
    categoryService = category_service.NewCategoryService(categoryRepo) // Adjust constructor as needed
}

func DeleteCategoryHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    idStr := event.PathParameters["id"]
    categoryID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil || categoryID <= 0 {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid category ID"), nil
    }

    if err := categoryService.DeleteCategory(categoryID); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Failed to delete category"), nil
    }

    return response.LambdaResponse(http.StatusOK, map[string]any{
        "status":  true,
        "message": "Category deleted successfully",
    }, ""), nil
}

func main() {
    lambda.Start(DeleteCategoryHandler)
}