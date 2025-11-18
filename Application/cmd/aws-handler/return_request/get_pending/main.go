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
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/return_request_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var returnRequestService return_request_service.ReturnRequestServiceInterface
var categoryRepo category_repo.CategoryRepo
var userRepo user_repo.UserRepo

func init() {
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic("Failed to connect to DynamoDB: " + err.Error())
    }
    
    categoryRepo = category_repo.NewCategoryDBRepo(dynamo)
    lenderRepo := lender_repo.NewLenderDBRepo(dynamo)
    userRepo = user_repo.NewUserDBRepo(dynamo,lenderRepo)
    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo,userRepo)
    orderRepo := order_repo.NewOrderDBRepo(dynamo,productRepo)
    rrRepo := return_request_repo.NewReturnRequestDBRepo(dynamo)

    returnRequestService = return_request_service.NewReturnRequestService(orderRepo, productRepo, rrRepo)
}


func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
    requests, err := returnRequestService.GetPendingReturnRequests(userCtx.ID)
    if err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, err.Error()), nil
    }

    body, _ := json.Marshal(map[string]interface{}{
        "status":   true,
        "requests": requests,
    })

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": "application/json"},
        Body:       string(body),
    }, nil
}

func main() {
    lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}