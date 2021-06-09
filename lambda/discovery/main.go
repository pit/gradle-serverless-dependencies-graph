package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"net/http"
	"terraform-serverless-private-registry/lib/helpers"
)

var (
	logger *zap.Logger
)

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
}

func main() {
	lambda.Start(Handler)
}

type Response struct {
	Modules   string `json:"modules.v1"`
	Providers string `json:"providers.v1"`
	//Login string `json:"modules.v1"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()
	logger.Debug("discovery called",
		zap.String("reqId", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	resp := new(Response)
	resp.Modules = fmt.Sprintf("https://%s/modules/v1", request.RequestContext.DomainName)
	resp.Providers = fmt.Sprintf("https://%s/providers/v1", request.RequestContext.DomainName)

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
