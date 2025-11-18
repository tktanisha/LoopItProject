package main

import (
	"context"
	"encoding/json"
	"net/http"

	"loopit/internal/api/middleware"
	"loopit/internal/db"
	"loopit/internal/enums"
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


    lenderRepo = lender_repo.NewLenderDBRepo(dynamo)
    userRepo := user_repo.NewUserDBRepo(dynamo,lenderRepo)
    categoryRepo = category_repo.NewCategoryDBRepo(dynamo)
    productRepo := product_repo.NewProductDBRepo(dynamo,categoryRepo,userRepo)
    orderRepo := order_repo.NewOrderDBRepo(dynamo,productRepo)
    returnRepo := return_request_repo.NewReturnRequestDBRepo(dynamo)

    orderService = order_service.NewOrderService(orderRepo, returnRepo, productRepo)
    productService = product_service.NewProductService(productRepo, userRepo)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest, userCtx *models.UserContext) (events.APIGatewayProxyResponse, error) {
    if userCtx.Role != enums.RoleLender {
        return response.LambdaResponse(http.StatusForbidden, nil, "Only lenders can view their orders"), nil
    }

    orders, err := orderService.GetLenderOrders(userCtx)
    if err != nil {
        return response.LambdaResponse(http.StatusInternalServerError, nil, err.Error()), nil
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