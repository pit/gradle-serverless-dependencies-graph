package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"gradle-serverless-dependencies-graph/lib/helpers"
	"gradle-serverless-dependencies-graph/lib/storage"
	"net/http"
	"os"
)

var (
	storageSvc *storage.Storage
	logger     *zap.Logger
)

type Response struct {
	Status       string  `json:"status"`
	UsedCapacity float64 `json:"used-capacity"`
}

func init() {
	storageTableName := os.Getenv("DYNAMODB_TABLE_STORAGE")
	dependenciesTableName := os.Getenv("DYNAMODB_TABLE_DEPENDENCIES")
	repositoriesTableName := os.Getenv("DYNAMODB_TABLE_REPOSITORIES")
	cfg := storage.StorageConfig{
		StorageTableName: &storageTableName,
		DependenciesTableName: &dependenciesTableName,
		RepositoriesTableName: &repositoriesTableName,
	}

	logger, _ = helpers.InitLogger("DEBUG", true)

	storageSvc, _ = storage.NewStorage(cfg, logger)
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

	repo := fmt.Sprintf("%s/%s", request.PathParameters["org"], request.PathParameters["repo"])
	ref := request.PathParameters["ref"]
	var deps storage.DependenciesRest
	json.Unmarshal([]byte(request.Body), &deps)

	logger.Debug("Request data",
		zap.String("repo", repo),
		zap.String("branch", ref),
		zap.Reflect("deps", deps),
	)

	resp, err := storageSvc.UpsertRepositoryInfo(request.RequestContext.RequestID, repo, ref, deps)

	if err != nil {
		return helpers.ApiErrorUnknown(), nil
	} else {
		return helpers.ApiResponse(http.StatusOK, Response{Status: "ok", UsedCapacity: resp.UsedCapacity}), nil
	}
}
