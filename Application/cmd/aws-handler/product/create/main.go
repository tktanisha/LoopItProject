package main

import (
	"context"
	"encoding/json"
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
var lenderRepo lender_repo.LenderRepo
	

func init() {
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic("Failed to connect to DynamoDB: " + err.Error())
    }
    lenderRepo = lender_repo.NewLenderDBRepo(dynamo) 
    userRepo := user_repo.NewUserDBRepo(dynamo,lenderRepo)
    categoryRepo = category_repo.NewCategoryDBRepo(dynamo)
    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo ,userRepo)
    productService = product_service.NewProductService(productRepo, userRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest,userCtx *models.UserContext ) (events.APIGatewayProxyResponse, error) {
    var product models.Product
    if err := json.Unmarshal([]byte(event.Body), &product); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    if err := productService.CreateProduct(&product, userCtx); err != nil {
        return response.LambdaResponse(http.StatusForbidden, nil, err.Error()), nil
    }

    return response.LambdaResponse(http.StatusCreated, map[string]interface{}{
        "status":  true,
        "message": "Product created successfully",
        "product": product,
    }, ""), nil
}

func main() {
    lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}