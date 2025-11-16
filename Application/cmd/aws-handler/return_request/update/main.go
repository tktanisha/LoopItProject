package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/enums/return_request_status"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
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

    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo,userRepo)
    orderRepo := order_repo.NewOrderDBRepo(dynamo,productRepo)
    rrRepo := return_request_repo.NewReturnRequestDBRepo(dynamo,nil)

    returnRequestService = return_request_service.NewReturnRequestService(orderRepo, productRepo, rrRepo, nil)
}


func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
    reqIDStr := event.PathParameters["requestId"]
    var reqID int64
    if _, err := fmt.Sscanf(reqIDStr, "%d", &reqID); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request ID"), nil
    }

    var payload struct {
        Status string `json:"status"`
    }
    if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    newStatus, err := return_request_status.ParseStatus(payload.Status)
    if err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid status value"), nil
    }

    if err := returnRequestService.UpdateReturnRequestStatus(userCtx.ID, reqID, newStatus); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

    return response.LambdaResponse(http.StatusOK, map[string]interface{}{
        "status":  true,
        "message": "Return request status updated successfully",
    }, ""), nil
}

func main() {
    lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}