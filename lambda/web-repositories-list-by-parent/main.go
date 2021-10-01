package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/smithy-go/ptr"
	"go.uber.org/zap"
	"gradle-serverless-dependencies-graph/lib/helpers"

	"gradle-serverless-dependencies-graph/lib/storage"
	"net/http"
	"os"
	"text/template"
)

var (
	storageSvc *storage.Storage
	logger     *zap.Logger
)

func init() {
	storageTableName := os.Getenv("DYNAMODB_TABLE_STORAGE")
	dependenciesTableName := os.Getenv("DYNAMODB_TABLE_DEPENDENCIES")
	repositoriesTableName := os.Getenv("DYNAMODB_TABLE_REPOSITORIES")
	cfg := storage.StorageConfig{
		StorageTableName:      &storageTableName,
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

	reqId := request.RequestContext.RequestID

	logger.Info("lambda called",
		zap.String("reqId", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	var parent *string
	org, okOrg := request.PathParameters["org"]
	repo, okRepo := request.PathParameters["repo"]

	if okOrg && okRepo {
		parent = ptr.String(fmt.Sprintf("%s/%s", org, repo))
	} else {
		parent = nil
	}

	resp, err := storageSvc.ListRepositoriesByParent(reqId, parent)

	if err != nil {
		return nil,err
	}

	data := struct {
		Items []storage.RepositoryDto
		Parent string
	}{
		Items: *resp,
		Parent: ptr.ToString(parent),
	}

	var tpl *template.Template
	var errTpl error
	if tpl, errTpl = template.New("html").Parse(Template); errTpl != nil {
		return nil, errTpl
	}

	var contentIO bytes.Buffer
	if errTpl = tpl.Execute(&contentIO, data); errTpl != nil {
		return nil, errTpl
	}

	content := contentIO.String()

	return helpers.HtmlResponse(http.StatusOK, &content), nil
}
