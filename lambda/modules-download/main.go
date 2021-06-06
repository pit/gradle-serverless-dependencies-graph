package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/modules"
	"terraform-serverless-private-registry/lib/storage"

	"net/http"
)

var (
	modulesSvc *modules.Modules
	logger     *zap.Logger
)

func init() {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ = helpers.InitLogger("DEBUG", true)
	logger.Debug("Lambda loading")

	storage, _ := storage.NewStorage(bucketName, logger)
	modulesSvc, _ = modules.NewModules(storage, logger)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger.Debug("Lambda called",
		zap.String("requestId", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	namespace := request.PathParameters["namespace"]
	name := request.PathParameters["name"]
	provider := request.PathParameters["provider"]
	version := request.PathParameters["version"]

	params := modules.InputParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
		Version:   &version,
	}
	resp, _ := modulesSvc.GetDownloadUrl(request.RequestContext.RequestID, params)

	lambdaResp, _ := helpers.ApiResponse(http.StatusNoContent, resp)
	lambdaResp.Headers["X-Terraform-Get"] = *resp

	return lambdaResp, nil
}
