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
    var payload struct {
        OrderID int64 `json:"order_id"`
    }
    if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    if err := returnRequestService.CreateReturnRequest(userCtx.ID, payload.OrderID); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

    return response.LambdaResponse(http.StatusCreated, map[string]interface{}{
        "status":  true,
        "message": "Return request created successfully",
    }, ""), nil
}

func main() {
    lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}