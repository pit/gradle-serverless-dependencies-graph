package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"os"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/modules"
	"terraform-serverless-private-registry/lib/storage"
)

var (
	modulesSvc *modules.Modules
	logger     *zap.Logger
)

func init() {
	bucketName := os.Getenv("BUCKET_NAME")
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
	resp, err := modulesSvc.GetDownloadUrl(request.RequestContext.RequestID, params)

	if err != nil {
		if err.Code == modules.ErrNotFound {
			return helpers.ApiNotFound(), nil
		}
	}

	lambdaResp := helpers.ApiNoContent()
	lambdaResp.Headers["X-Terraform-Get"] = *resp

	return lambdaResp, nil
}
