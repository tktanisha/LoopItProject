package main

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"

// 	"loopit/internal/db"
// 	"loopit/internal/enums"
// 	"loopit/internal/models"
// 	"loopit/internal/repository/lender_repo"
// 	"loopit/internal/repository/user_repo"
// 	"loopit/internal/services/user_service"
// 	response "loopit/internal/utils/lambda_response"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// )

// var userService user_service.UserServiceInterface
// var lenderRepo lender_repo.LenderRepo

// func init() {
//     dynamo, err := db.ConnectDynamo()
//     if err != nil {
//         panic("Failed to connect to DynamoDB: " + err.Error())
//     }
//     userRepo := user_repo.NewUserDBRepo(dynamo,lenderRepo)
//     userService = user_service.NewUserService(userRepo)
// }

// func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//     var userCtx models.UserContext
//     if err := json.Unmarshal([]byte(event.Body), &userCtx); err != nil {
//         return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid user context"), nil
//     }

//     if err := userService.BecomeLender(&userCtx); err != nil {
//         return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
//     }

//     userCtx.Role = enums.RoleLender
//     return response.LambdaResponse(http.StatusOK, map[string]interface{}{
//         "status":  true,
//         "message": "User promoted to lender successfully",
//         "user":    userCtx,
//     }, ""), nil
// }

// func main() {
//     lambda.Start(Handler)
// }