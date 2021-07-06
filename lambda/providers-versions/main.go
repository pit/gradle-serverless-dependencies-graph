package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"net/http"
	"os"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/providers"
	"terraform-serverless-private-registry/lib/storage"
)

var (
	providersSvc *providers.Providers
	logger       *zap.Logger
)

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
	bucketName := os.Getenv("BUCKET_NAME")
	storage, _ := storage.NewStorage(bucketName, logger)
	providersSvc, _ = providers.NewProviders(storage, logger)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()

	//requestJson, _ := json.Marshal(request)
	reqId := request.RequestContext.RequestID

	logger.Info("lambda called",
		zap.String("reqId", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	providerNamespace := request.PathParameters["namespace"]
	providerType := request.PathParameters["type"]

	params := providers.ListProviderVersionsInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
	}
	resp, err := providersSvc.ListProviderVersions(reqId, params)

	if err != nil {
		if err.Code == providers.ErrNotFound {
			helpers.ApiErrorNotFound()
		}
	}

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
