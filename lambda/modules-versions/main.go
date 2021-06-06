package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"net/http"
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
	logger, _ = helpers.InitLogger("DEBUG", true)
	bucketName := os.Getenv("BUCKET_NAME")
	storage, _ := storage.NewStorage(bucketName, logger)
	modulesSvc, _ = modules.NewModules(storage, logger)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()

	requestJson, _ := json.Marshal(request)
	reqId := request.RequestContext.RequestID

	logger.Info(fmt.Sprintf("%s lambda called", reqId),
		zap.ByteString("request", requestJson),
	)

	namespace := request.PathParameters["namespace"]
	name := request.PathParameters["name"]
	provider := request.PathParameters["provider"]

	params := modules.InputParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
	}
	resp, err := modulesSvc.ListModuleVersions(reqId, params)

	if err != nil {
		if err.Code == modules.ErrNotFound {
			helpers.ApiNotFound()
		}
	}

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
