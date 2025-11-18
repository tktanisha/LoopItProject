package main

import (
	"context"
	"net/http"

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
	filters := models.UserFilter{
		Search:    event.QueryStringParameters["search"],
		Role:      event.QueryStringParameters["role"],
		SocietyID: event.QueryStringParameters["society_id"],
	}

	users, err := userService.GetAllUsers(filters)
	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, nil, err.Error()), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]interface{}{
		"status": true,
		"users":  users,
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}
