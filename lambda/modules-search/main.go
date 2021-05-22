package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            request.Body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}
