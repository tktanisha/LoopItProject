package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/repository/society_repo"
	"loopit/internal/services/society_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)



 var societyService society_service.SocietyServiceInterface
 var societyRepo society_repo.SocietyRepo

func init() {
    
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to DynamoDB: %v", err))
    }
    log.Print("entered in the init ")

    societyRepo = society_repo.NewSocietyDBRepo(dynamo)
    societyService = society_service.NewSocietyService(societyRepo)
}
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    societies, err := societyService.GetAllSocieties()
    if err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, "Failed to fetch societies"), nil
    }
    return response.LambdaResponse(http.StatusOK, map[string]interface{}{
        "status":    true,
        "societies": societies,
    }, ""), nil
}
func main() {
    lambda.Start(middleware.WithCORS(Handler))
}