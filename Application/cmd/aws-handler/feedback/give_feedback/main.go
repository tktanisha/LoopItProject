package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/feedback_repo"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/feedback_service"
	"loopit/internal/utils"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var feedbackService feedback_service.FeedbackServiceInterface
var categoryRepo category_repo.CategoryRepo
var userRepo user_repo.UserRepo

func init() {
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic("Failed to connect to DynamoDB: " + err.Error())
    }
    
    lenderRepo:= lender_repo.NewLenderDBRepo(dynamo)
    userRepo = user_repo.NewUserDBRepo(dynamo,lenderRepo)
    categoryRepo = category_repo.NewCategoryDBRepo(dynamo)
    feedbackRepo := feedback_repo.NewFeedBackDBRepo(dynamo)
    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo,userRepo)
    orderRepo := order_repo.NewOrderDBRepo(dynamo,productRepo)

    feedbackService = feedback_service.NewFeedbackService(feedbackRepo, productRepo, orderRepo)
}
func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
    var payload struct {
        OrderID      string `json:"order_id"`
        FeedbackText string `json:"feedback_text"`
        Rating       int    `json:"rating"`
    }
    if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    orderId, errs := strconv.ParseInt(payload.OrderID, 10, 64)
    if errs != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid order ID format"), nil
    }

    // Basic validations
    if err := utils.ValidateOrderID(orderId); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }
    if err := utils.ValidateFeedbackText(payload.FeedbackText); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }
    if err := utils.ValidateRating(payload.Rating); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

    if err := feedbackService.GiveFeedback(orderId, payload.FeedbackText, payload.Rating, userCtx); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

    return response.LambdaResponse(http.StatusCreated, map[string]interface{}{
        "status":  true,
        "message": "Feedback given successfully",
    }, ""), nil
}

func main() {
    lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}