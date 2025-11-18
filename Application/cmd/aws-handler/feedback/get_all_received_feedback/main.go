package main

import (
	"context"
	"encoding/json"
	"net/http"

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
    feedbacks, err := feedbackService.GetAllReceivedFeedbacks(userCtx)
    if err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, "Failed to fetch received feedbacks"), nil
    }

    body, _ := json.Marshal(map[string]interface{}{
        "status":    true,
        "feedbacks": feedbacks,
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