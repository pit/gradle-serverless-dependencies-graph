package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"terraform-serverless-private-registry/lib"
	"net/http"
)

func main() {
	lambda.Start(Handler)
}

type Response struct {
	Modules string `json:"modules.v1"`
	Providers string `json:"providers.v1"`
	//Login string `json:"modules.v1"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctxS,_ := json.Marshal(ctx)
	reqS,_ := json.Marshal(request)
	fmt.Printf("ctx: %s", ctxS)
	fmt.Printf("request: %s", reqS)
	resp := new(Response)
	resp.Modules = fmt.Sprintf("https://%s/modules/v1", request.RequestContext.DomainName)
	resp.Providers = fmt.Sprintf("https://%s/providers/v1", request.RequestContext.DomainName)
	return lib.ApiResponse(http.StatusOK, resp)
}
