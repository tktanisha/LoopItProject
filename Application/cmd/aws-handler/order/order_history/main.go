package main

import (
	"context"
	"encoding/json"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/enums/order_status"
	"loopit/internal/models"
	"loopit/internal/repository/category_repo"
	"loopit/internal/repository/lender_repo"
	"loopit/internal/repository/order_repo"
	"loopit/internal/repository/product_repo"
	"loopit/internal/repository/return_request_repo"
	"loopit/internal/repository/user_repo"
	"loopit/internal/services/order_service"
	"loopit/internal/services/product_service"
	response "loopit/internal/utils/lambda_response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)
var orderService order_service.OrderServiceInterface
var productService product_service.ProductServiceInterface
var lenderRepo lender_repo.LenderRepo
var categoryRepo category_repo.CategoryRepo

func init() {
    dynamo, err := db.ConnectDynamo()
    if err != nil {
        panic("Failed to connect to DynamoDB: " + err.Error())
    }

    userRepo := user_repo.NewUserDBRepo(dynamo,lenderRepo)
    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo,userRepo)
    orderRepo := order_repo.NewOrderDBRepo(dynamo,productRepo)
    returnRepo := return_request_repo.NewReturnRequestDBRepo(dynamo)

    orderService = order_service.NewOrderService(orderRepo, returnRepo, productRepo)
    productService = product_service.NewProductService(productRepo, userRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
    statusStr := event.QueryStringParameters["status"]
    var filterStatus []order_status.Status
    if statusStr != "" {
        st, err := order_status.ParseStatus(statusStr)
        if err != nil {
            return response.LambdaResponse(http.StatusBadRequest, nil, "Invalid status filter"), nil
        }
        filterStatus = append(filterStatus, st)
    }

    orders, err := orderService.GetOrderHistory(userCtx, filterStatus)
    if err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, "Failed to fetch order history"), nil
    }

    var orderResponses []*models.OrderDto
    for _, order := range orders {
        product, err := productService.GetProductByID(order.ProductID)
        if err != nil {
            continue
        }
        orderResponses = append(orderResponses, &models.OrderDto{Order: *order, Product: *product})
    }

    body, _ := json.Marshal(map[string]interface{}{
        "status": true,
        "orders": orderResponses,
    })

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": "application/json"},
        Body:       string(body),
    }, nil
}

func main() {
    lambda.Start(middleware.WithCORS(middleware.WithAuth(Handler)))
}