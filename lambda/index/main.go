package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gradle-serverless-dependencies-graph/lib/helpers"
	"net/http"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	resp := `
<html><body><pre>
<a href="/dependency/">Dependencies</a>
<a href="/repositories/">Repositories</a>
</pre></body></html>
`
	return helpers.HtmlResponse(http.StatusOK, &resp), nil
}
