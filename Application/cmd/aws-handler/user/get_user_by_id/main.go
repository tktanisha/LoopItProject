package main

import (
	"context"
	"net/http"
	"strconv"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/user_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var userService user_service.UserServiceInterface
var lenderRepo lender_repo.LenderRepo

func init() {
	dynamo, err := db.ConnectDynamo()
	if err != nil {
		panic("Failed to connect to DynamoDB: " + err.Error())
	}
	lenderRepo = lender_repo.NewLenderDBRepo(dynamo)
	userRepo := user_repo.NewUserDBRepo(dynamo, lenderRepo)
	userService = user_service.NewUserService(userRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
	idStr := event.PathParameters["id"]
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || userID <= 0 {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid user ID"), nil
	}

	user, err := userService.GetUserByID(userID)
	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, nil, err.Error()), nil
	}
	if user == nil {
		return response.LambdaResponse(http.StatusNotFound, nil, "User not found"), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]interface{}{
		"status": true,
		"user":   user,
		"message": "successfully fetched user",
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}
