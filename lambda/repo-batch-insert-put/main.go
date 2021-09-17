package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"gradle-serverless-dependencies-graph/lib/helpers"
)

var (
	//modulesSvc *modules.Modules
	logger     *zap.Logger
)

func init() {
	//dynamodbTable := os.Getenv("DYNAMODB_TABLE")
	logger, _ = helpers.InitLogger("DEBUG", true)
	logger.Debug("Lambda loading")

	//storage, _ := storage.NewStorage(dynamodbTable, logger)
	//modulesSvc, _ = modules.NewModules(storage, logger)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()
	logger.Debug("Lambda called",
		zap.String("requestId", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	repo := request.PathParameters["repo"]
	branch := request.PathParameters["branch"]

	logger.Debug("Request data",
		zap.String("repo", repo),
		zap.String("branch", branch),
	)

	return helpers.ApiErrorNotFound(), nil
	//if err != nil {
	//	if err.Code == modules.ErrNotFound {
	//		return helpers.ApiErrorNotFound(), nil
	//	}
	//}

	//lambdaResp := helpers.ApiErrorNoContent()
	//lambdaResp.Headers["X-Terraform-Get"] = *resp

	//return lambdaResp, nil
}
