package main

import (
	"context"
	"net/http"
	"strconv"

	"loopit/internal/db"
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
    userRepo := user_repo.NewUserDBRepo(dynamo,lenderRepo)
    userService = user_service.NewUserService(userRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    idStr := event.PathParameters["id"]
    userID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil || userID <= 0 {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid user ID"), nil
    }

    if err := userService.DeleteUserByID(userID); err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, "Failed to delete user"), nil
    }

    return response.LambdaResponse(http.StatusOK, map[string]interface{}{
        "status":  true,
        "message": "User deleted successfully",
    }, ""), nil
}
func main() {
    lambda.Start(Handler)
}