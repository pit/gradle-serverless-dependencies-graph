package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"terraform-serverless-private-registry/lib/helpers"
)

func main() {
	spew.Config.Indent = "  "
	spew.Config.DisableMethods = true
	spew.Config.DisablePointerMethods = true

	lambda.Start(Handler)
}

type Response struct {
	Modules   string `json:"modules.v1"`
	Providers string `json:"providers.v1"`
	//Login string `json:"modules.v1"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("ctx: %s, request: %s", spew.Sdump(ctx), spew.Sdump(request))

	resp := new(Response)
	resp.Modules = fmt.Sprintf("https://%s/modules/v1", request.RequestContext.DomainName)
	resp.Providers = fmt.Sprintf("https://%s/providers/v1", request.RequestContext.DomainName)

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
