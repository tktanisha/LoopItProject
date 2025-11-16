package main

import (
	"context"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/product_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var productService product_service.ProductServiceInterface
var categoryRepo category_repo.CategoryRepo
var LenderRepo lender_repo.LenderRepo
	

func init() {
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic("Failed to connect to DynamoDB: " + err.Error())
    }
    userRepo := user_repo.NewUserDBRepo(dynamo,LenderRepo)
    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo ,userRepo)
    productService = product_service.NewProductService(productRepo, userRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    filters := models.ProductFilter{
        Search:      event.QueryStringParameters["search"],
        LenderID:    event.QueryStringParameters["lender_id"],
        CategoryID:  event.QueryStringParameters["category_id"],
        IsAvailable: event.QueryStringParameters["is_available"],
    }

    products, err := productService.GetAllProducts(filters)
    if err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, "Failed to fetch products"), nil
    }

    return response.LambdaResponse(http.StatusOK, map[string]interface{}{
        "status":   true,
        "products": products,
    }, ""), nil
}

func main() {
  lambda.Start(middleware.WithCORS(Handler))}