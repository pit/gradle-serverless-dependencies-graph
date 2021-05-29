package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"net/http"
	"terraform-serverless-private-registry/lib"
)

var (
	modules *lib.Modules
	logger *zap.Logger
)

func init(){
	bucketName := "terraform-registry-kvinta-io"
	logger,_ := lib.InitLogger("DEBUG", true)
	storage,_ := lib.NewStorage(&bucketName, logger)
	modules,_ = lib.NewModules(storage, logger)
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

	resp,_ := modules.ListModuleVersions(ctx, namespace, name, provider)

	return lib.ApiResponse(http.StatusOK, resp)
}
