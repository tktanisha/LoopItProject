package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
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

func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
	idStr := event.PathParameters["id"]
	societyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || societyID <= 0 {
		log.Print("invaid society id")
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid society ID"), nil
	}
	log.Print("enterd in the handler")

	var payload struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Pincode  string `json:"pincode"`
	}
	if err := json.Unmarshal([]byte(event.Body), &payload); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid request payload"), nil
	}

	if err := societyService.UpdateSociety(societyID, payload.Name, payload.Location, payload.Pincode); err != nil {
		log.Print("after updating ")
		return response.LambdaResponse(http.StatusInternalServerError, nil, err.Error()), nil
	}
    
	return response.LambdaResponse(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Society updated successfully",
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}
