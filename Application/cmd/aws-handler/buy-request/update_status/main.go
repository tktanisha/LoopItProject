package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/enums/buyer_request_status"
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
	lenderRepo = lender_repo.NewLenderDBRepo(dynamo)
	userRepo := user_repo.NewUserDBRepo(dynamo, lenderRepo)
	productRepo := product_repo.NewProductDBRepo(dynamo, categoryRepo, userRepo)
	orderRepo := order_repo.NewOrderDBRepo(dynamo, productRepo)

	buyerRequestService = buyer_request_service.NewBuyerRequestService(buyerReqRepo, productRepo, orderRepo, categoryRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
	reqIDStr := event.PathParameters["requestId"]
	reqID, err := strconv.ParseInt(reqIDStr, 10, 64)
	log.Print(reqID)
	if err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid buyer request ID"), nil
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
	}

	statusEnum, err := buyer_request_status.ParseStatus(payload.Status)
	if err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid status value"), nil
	}

	if err := buyerRequestService.UpdateBuyerRequestStatus(reqID, statusEnum, userCtx); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Buyer request status updated successfully",
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}
