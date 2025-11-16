package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/repository/society_repo"
	"loopit/internal/services/society_service"
	"loopit/internal/utils"
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
func CreateSocietyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var payload struct {
        Name     string `json:"name"`
        Location string `json:"location"`
        Pincode  string `json:"pincode"`
    }
    if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
    }

    if err := utils.ValidateSociety(payload.Name, payload.Location, payload.Pincode); err != nil {
        return response.LambdaResponse(http.StatusBadRequest, nil, err.Error()), nil
    }

    if err := societyService.CreateSociety(payload.Name, payload.Location, payload.Pincode); err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, "Failed to create society"), nil
    }

    return response.LambdaResponse(http.StatusCreated, map[string]interface{}{
        "status":  true,
        "message": "Society created successfully",
    }, ""), nil
}
func main() {
    lambda.Start(middleware.WithCORS(CreateSocietyHandler))

}