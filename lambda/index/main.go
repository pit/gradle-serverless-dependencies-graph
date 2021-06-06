package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"terraform-serverless-private-registry/lib/helpers"
)

func main() {
	lambda.Start(Handler)
}

type Response struct {
	Status string `json:"status"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	resp := new(Response)
	resp.Status = "OK"
	return helpers.ApiResponse(http.StatusOK, resp), nil
}
