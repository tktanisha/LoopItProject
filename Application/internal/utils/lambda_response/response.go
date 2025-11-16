package response

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func LambdaResponse(status int, data any, message string) events.APIGatewayProxyResponse {
    body, _ := json.Marshal(map[string]any{
        "message": message,
        "data":    data,
    })

    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Headers: map[string]string{
            "Content-Type":"application/json",
        },
        Body: string(body),
    }
}

