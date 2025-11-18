package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

type registerRequest struct {
    FullName    string `json:"fullName"`
    Email       string `json:"email"`
    Password    string `json:"password"`
    PhoneNumber string `json:"phoneNumber"`
    Address     string `json:"address"`
    SocietyID   int64  `json:"societyId"`
}

var authService auth_service.AuthServiceInterface
var userRepo user_repo.UserRepo
var lenderRepo lender_repo.LenderRepo

func init() {
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to DynamoDB: %v", err))
    }
    log.Print("Init: DynamoDB connected")

    userRepo = user_repo.NewUserDBRepo(dynamo, lenderRepo)
    authService = auth_service.NewAuthService(userRepo)
}

func RegisterHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var req registerRequest
    if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    if err := utils.ValidateFullName(req.FullName); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }
    if err := utils.ValidateEmail(req.Email); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }
    if err := utils.ValidatePassword(req.Password); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }
    if err := utils.ValidatePhoneNumber(req.PhoneNumber); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }
    if err := utils.ValidateAddress(req.Address); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

    // Register user
    err := authService.Register(&models.User{
        FullName:     req.FullName,
        Email:        req.Email,
        PasswordHash: req.Password,
        PhoneNumber:  req.PhoneNumber,
        Address:      req.Address,
        SocietyID:    req.SocietyID,
        CreatedAt:    time.Now(),
    })
    if err != nil {
        log.Printf("Registration failed: %v", err)
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

	log.Print("before re login----------")

    // Auto-login after registration
    token, user, err := authService.Login(req.Email, req.Password)
    if err != nil {
        log.Printf("Auto login failed: %v", err)
        return response.LambdaResponse(http.StatusUnauthorized, nil, "Invalid credentials"), nil
    }

	log.Print("after re login====",user,token)
	

    userCtx := &models.UserContext{
        ID:   user.ID,
        Name: user.FullName,
        Role: user.Role,
    }

    return response.LambdaResponse(http.StatusOK, map[string]any{
        "token": token,
        "user":  userCtx,
    }, "Registration successful"), nil
}

func main() {
	log.Print("entered in lambda")
    lambda.Start(middleware.WithCORS(RegisterHandler))
}