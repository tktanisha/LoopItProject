package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/internal/repository/buyer_request_repo"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/buyer_request_service"
	"loopit/internal/services/product_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var buyerRequestService buyer_request_service.BuyerRequestServiceInterface
var lenderRepo lender_repo.LenderRepo
var productService product_service.ProductServiceInterface

func init() {
	dynamo, err := db.ConnectDynamo()
	if err != nil {
		panic("Failed to connect to DynamoDB: " + err.Error())
	}

	buyerReqRepo := buyer_request_repo.NewBuyerRequestDBRepo(dynamo)
	categoryRepo := category_repo.NewCategoryDBRepo(dynamo)
	lenderRepo = lender_repo.NewLenderDBRepo(dynamo)
	userRepo := user_repo.NewUserDBRepo(dynamo, lenderRepo)
	productRepo := product_repo.NewProductDBRepo(dynamo, categoryRepo, userRepo)
	orderRepo := order_repo.NewOrderDBRepo(dynamo, productRepo)

	buyerRequestService = buyer_request_service.NewBuyerRequestService(buyerReqRepo, productRepo, orderRepo, categoryRepo)
	productService = product_service.NewProductService(productRepo, userRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
	var productID *int64
	if productIDStr := event.QueryStringParameters["product_id"]; productIDStr != "" {
		id, err := strconv.ParseInt(productIDStr, 10, 64)
		if err != nil {
			return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid product_id query param"), nil
		}
		productID = &id
	}

	var statusFilter []string
	if statusStr := event.QueryStringParameters["status"]; statusStr != "" {
		statusFilter = strings.Split(statusStr, ",")
	}

	requests, err := buyerRequestService.GetAllBuyerRequests(productID, statusFilter)
	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, nil, err.Error()), nil
	}

	var requestResponses []models.BuyingRequestDto
	for _, req := range requests {
		product, err := productService.GetProductByID(req.ProductID)
		if err != nil {
			continue
		}
		requestResponses = append(requestResponses, models.BuyingRequestDto{
			BuyRequest: req,
			Product:    *product,
		})
	}

	return response.LambdaResponse(http.StatusOK, map[string]interface{}{
		"status":   true,
		"requests": requestResponses,
	}, ""), nil
}

func main() {
	lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}
