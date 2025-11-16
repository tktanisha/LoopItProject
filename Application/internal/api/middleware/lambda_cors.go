package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

var defaultHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Headers": "Content-Type,Authorization",
	"Access-Control-Allow-Methods": "OPTIONS,GET,POST,PUT,PATCH,DELETE",
}

func WithCORS(
	fn func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error),
) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if req.HTTPMethod == "OPTIONS" {
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers:    cloneHeaders(nil),
			}, nil
		}

		res, err := fn(ctx, req)
		res.Headers = cloneHeaders(res.Headers)
		return res, err
	}
}

func cloneHeaders(existing map[string]string) map[string]string {
	headers := make(map[string]string, len(defaultHeaders))
	for k, v := range defaultHeaders {
		headers[k] = v
	}

	for k, v := range existing {
		headers[k] = v
	}

	return headers
}
