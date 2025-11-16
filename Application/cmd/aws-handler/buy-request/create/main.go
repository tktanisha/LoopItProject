package main

import (
	"context"
	"encoding/json"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/buyer_request_repo"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/buyer_request_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var buyerRequestService buyer_request_service.BuyerRequestServiceInterface
var lenderRepo lender_repo.LenderRepo

func init() {
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic("Failed to connect to DynamoDB: " + err.Error())
    }

    buyerReqRepo := buyer_request_repo.NewBuyerRequestDBRepo(dynamo)
    categoryRepo := category_repo.NewCategoryDBRepo(dynamo)
    userRepo := user_repo.NewUserDBRepo(dynamo,lenderRepo)
    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo,userRepo)
    orderRepo := order_repo.NewOrderDBRepo(dynamo,productRepo)

    buyerRequestService = buyer_request_service.NewBuyerRequestService(buyerReqRepo, productRepo, orderRepo, categoryRepo, nil)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
    var payload struct {
        ProductID int64 `json:"product_id"`
    }
    if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    if err := buyerRequestService.CreateBuyerRequest(payload.ProductID, userCtx); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

    return response.LambdaResponse(http.StatusCreated, map[string]interface{}{
        "status":  true,
        "message": "Buyer request created successfully",
    }, ""), nil
}

func main() {
    lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}