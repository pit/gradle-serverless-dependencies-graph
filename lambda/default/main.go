package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gradle-serverless-dependencies-graph/lib/helpers"
	"net/http"
)

func main() {
	lambda.Start(Handler)
}

type Response struct {
	Status string `json:"status"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	resp := new(Response)
	resp.Status = "What do you need?"
	return helpers.ApiResponse(http.StatusTeapot, resp), nil
}
