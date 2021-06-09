package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"os"
	"strings"
	"terraform-serverless-private-registry/lib/authorization"
	"terraform-serverless-private-registry/lib/helpers"
)

var (
	logger           *zap.Logger
	authorizationSvc *authorization.Authorization
)

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
	authorizationSvc, _ = authorization.NewAuthorization(logger)
}

func main() {
	lambda.Start(Handler)
}

func Handler(request APIGatewayAuthorizerRequest) (*APIGatewayAuthorizerResponse, error) {
	defer logger.Sync()

	logger.Debug("Running authorizer request",
		zap.String("reqId", request.RequestContext.RequestID),
		zap.Reflect("request.Type", request.Type),
		zap.Reflect("request.Version", request.Version),
		zap.Reflect("request.RouteArn", request.RouteArn),
		zap.Reflect("request.RouteKey", request.RouteKey),
		zap.Reflect("request.RequestContext", request.RequestContext),
		zap.Reflect("request.PathParameters", request.PathParameters),
		zap.Reflect("request.StageVariables", request.StageVariables),
		zap.Strings("request.IdentitySource.md5", helpers.GenerateMD5List(request.IdentitySource)),
	)

	region := os.Getenv("AWS_REGION")
	resourceArn := fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/%s/GET/*",
		region,
		request.RequestContext.AccountID,
		request.RequestContext.ApiId,
		request.RequestContext.Stage,
	)

	if len(request.IdentitySource) >= 1 {
		identitySource := strings.Split(request.IdentitySource[0], " ")
		logger.Debug("Authorization to check",
			zap.Strings("identitySource.IdentitySource.split.md5", helpers.GenerateMD5List(identitySource)),
		)
		if len(identitySource) == 2 {
			authEncodedBytes, err := base64.StdEncoding.DecodeString(identitySource[1])
			if err != nil {
				return nil, err
			}
			authEncoded := string(authEncodedBytes)
			logger.Debug("Authorization Info to check",
				zap.String("authEncoded.md5", helpers.GenerateMD5(authEncoded)),
			)

			authInfo := strings.Split(authEncoded, ":")
			if len(authInfo) == 2 && authorizationSvc.CheckCredentials(request.RequestContext.RequestID, authInfo[0], authInfo[1]) {
				resp := generatePolicy(authInfo[0],"Allow",resourceArn)
				logger.Debug("Response",
					zap.Reflect("resp", resp),
				)
				return resp, nil
			}
		}
	}
	resp := generatePolicy("", "Deny", resourceArn)
	logger.Debug("Response",
		zap.Reflect("resp", resp),
	)
	return resp, nil
}

func generatePolicy(principalId string, effect string, resourceArn string) *APIGatewayAuthorizerResponse {
	return &APIGatewayAuthorizerResponse{
		PrincipalID: principalId,
		PolicyDocument: IAMPolicyDocument{
			Version: "2012-10-17",
			Statement: []IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resourceArn},
				},
			},
		},
	}
}
