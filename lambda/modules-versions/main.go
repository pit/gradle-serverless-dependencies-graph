package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"net/http"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/modules"
	"terraform-serverless-private-registry/lib/storage"
)

var (
	modulesSvc *modules.Modules
	logger     *zap.Logger
)

const bucketName string = "terraform-registry-kvinta-io"

func init() {
	var err error
	logger, err = helpers.InitLogger("DEBUG", true)
	if err != nil {
		panic("Cannot start logger")
	} else {
		logger.Info("Starting lambda init")
	}
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

	}

	return helpers.ApiResponse(http.StatusOK, resp)
}
