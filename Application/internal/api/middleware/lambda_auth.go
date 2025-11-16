package middleware

import (
	"context"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/utils"
	response "loopit/internal/utils/lambda_response"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func WithAuth(
    fn func(context.Context, events.APIGatewayProxyRequest, *models.UserContext) (events.APIGatewayProxyResponse, error),
) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        authHeader := req.Headers["Authorization"]
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return response.LambdaResponse(http.StatusUnauthorized,nil, "missing or invalid authorization header"), nil
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := utils.ValidateJWT(token)
        if err != nil {
            return response.LambdaResponse(http.StatusUnauthorized, nil,"invalid or expired token"), nil
        }

        role, err := enums.ParseRole(claims.Role)
        if err != nil {
            return response.LambdaResponse(http.StatusUnauthorized, nil,"invalid role"), nil
        }

        userCtx := &models.UserContext{
            ID:   claims.UserID,
            Role: role,
        }

        return fn(ctx, req, userCtx)
    }
}

