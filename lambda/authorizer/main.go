package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"strings"
	"terraform-serverless-private-registry/lib/authorization"
	"terraform-serverless-private-registry/lib/helpers"
)

var (
	logger           *zap.Logger
	authorizationSvc *authorization.Authorization
)

type APIGatewayAuthorizerRequestContext struct {
	AccountID    string                                    `json:"accountId"`
	ApiId        string                                    `json:"apiId"`
	Path         string                                    `json:"path"`
	DomainName   string                                    `json:"domainName"`
	DomainPrefix string                                    `json:"domainPrefix"`
	Http         map[string]string                         `json:"http"`
	RequestID    string                                    `json:"requestId"`
	RouteKey     string                                    `json:"routeKey"`
	Stage        string                                    `json:"stage"`
	Identity     APIGatewayCustomAuthorizerRequestIdentity `json:"identity"`
	ResourcePath string                                    `json:"resourcePath"`
	// Authentication AuthenticationType `json:"authentication"`
}

type APIGatewayCustomAuthorizerRequestIdentity struct {
	APIKey   string `json:"apiKey"`
	SourceIP string `json:"sourceIp"`
}

type APIGatewayAuthorizerRequest struct {
	Version               string              `json:"version"`
	Type                  string              `json:"type"`
	RouteArn              string              `json:"routeArn"`
	RouteKey              string              `json:"routeKey"`
	IdentitySource        []string            `json:"identitySource"`
	RawQueryString        string              `json:"rawQueryString"`
	RawPath               string              `json:"rawPath"`
	Cookies               []string            `json:"cookies"`
	Headers               map[string]string   `json:"headers"`
	MultiValueHeaders     map[string][]string `json:"multiValueHeaders"`
	QueryStringParameters map[string]string   `json:"queryStringParameters"`
	//MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"` // TODO RnD this field via jsonRawMessage
	PathParameters map[string]string                  `json:"pathParameters"`
	StageVariables map[string]string                  `json:"stageVariables"`
	RequestContext APIGatewayAuthorizerRequestContext `json:"requestContext"`
}

//type APIGatewayAuthorizerResponse struct {
//	IsAuthorized bool `json"isAuthorized"`
//}
type APIGatewayAuthorizerResponse struct {
	PrincipalID        string                 `json:"principalId"`
	PolicyDocument     IAMPolicyDocument      `json:"policyDocument"`
	Context            map[string]interface{} `json:"context,omitempty"`
	UsageIdentifierKey string                 `json:"usageIdentifierKey,omitempty"`
}

// APIGatewayCustomAuthorizerPolicy represents an IAM policy
type IAMPolicyDocument struct {
	Version   string               `json:"Version"`
	Statement []IAMPolicyStatement `json:"Statement"`
}

// IAMPolicyStatement represents one statement from IAM policy with action, effect and resource
type IAMPolicyStatement struct {
	Action   []string `json:"Action"`
	Effect   string   `json:"Effect"`
	Resource []string `json:"Resource"`
}

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
	authorizationSvc, _ = authorization.NewAuthorization(logger)
}

func main() {
	lambda.Start(Handler)
}

//func Handler(ctx context.Context, request APIGatewayAuthorizerRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
func Handler(ctx context.Context, rawRequest map[string]json.RawMessage) (*APIGatewayAuthorizerResponse, error) {
	defer logger.Sync()
	logger.Debug("Running authorizer request",
		zap.Reflect("request", rawRequest),
	)

	var request *APIGatewayAuthorizerRequest
	rawJson, err := json.Marshal(rawRequest)
	if err != nil {
		logger.Error("Error with request json",
			zap.Error(err),
		)
	}

	err = json.Unmarshal(rawJson, &request)
	if err != nil {
		logger.Error("Error with request json",
			zap.Error(err),
		)
	}

	logger.Debug("Running authorizer request",
		zap.String("reqId", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	if len(request.IdentitySource) >= 1 {
		authorization := strings.Split(request.IdentitySource[0], " ")
		logger.Debug("Authorization to check",
			zap.Strings("authorization", authorization),
		)
		if len(authorization) == 2 {
			authEncodedBytes, err := base64.StdEncoding.DecodeString(authorization[1])
			if err != nil {
				return nil, err
			}
			authEncoded := string(authEncodedBytes)
			logger.Debug("Authorization Info to check",
				zap.String("authEncoded", authEncoded),
			)

			authInfo := strings.Split(authEncoded, ":")
			if len(authInfo) == 2 && authorizationSvc.CheckCredentials(request.RequestContext.RequestID, authInfo[0], authInfo[1]) {
				resp := generatePolicy(authInfo[0],"Allow","")
				logger.Debug("Response",
					zap.Reflect("resp", resp),
				)
				return resp, nil
				//return generatePolicy(authInfo[0], "Allow", request.RouteArn), nil
			}
		}
	}
	resp := generatePolicy("", "Deny", "")
	logger.Debug("Response",
		zap.Reflect("resp", resp),
	)
	return resp, nil
	//return generatePolicy("", "Deny", request.RouteArn), nil
}

//func generatePolicy(isAuthorized bool) *APIGatewayAuthorizerResponse {
//	return &APIGatewayAuthorizerResponse{
//		IsAuthorized: isAuthorized,
//	}
//}

func generatePolicy(principalId string, effect string, methodArn string) *APIGatewayAuthorizerResponse {
	return &APIGatewayAuthorizerResponse{
		PrincipalID: principalId,
		PolicyDocument: IAMPolicyDocument{
			Version: "2012-10-17",
			Statement: []IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{"arn:aws:execute-api:eu-central-1:463422107539:cdq0eq1dde/*"},
					//Resource: []string{methodArn},
				},
			},
		},
	}
}
