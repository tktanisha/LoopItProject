package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/auth_service"
	"loopit/internal/utils"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}


var authService auth_service.AuthServiceInterface
var userRepo user_repo.UserRepo
var lenderRepo lender_repo.LenderRepo



func init() {

    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to DynamoDB: %v", err))
    }
    log.Print("entered in the init ")

    userRepo = user_repo.NewUserDBRepo(dynamo, lenderRepo) 
    authService = auth_service.NewAuthService(userRepo)
    
}

func LoginHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var req loginRequest
    if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    if err := utils.ValidateEmail(req.Email); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid email"), nil
    }
    if err := utils.ValidatePassword(req.Password); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid password"), nil
    }

    token, user, err := authService.Login(req.Email, req.Password)
    if err != nil {
        log.Print(err)
        return response.LambdaResponse(http.StatusUnauthorized, nil, "Invalid credentials"), nil
    }

    userCtx := &models.UserContext{
        ID:   user.ID,
        Name: user.FullName,
        Role: user.Role,
    }


  
    

    return response.LambdaResponse(http.StatusOK, map[string]any{
        "token": token,
        "user":  userCtx,
    }, "Login successful"), nil
}

func main() {
    lambda.Start(middleware.WithCORS(LoginHandler))
}