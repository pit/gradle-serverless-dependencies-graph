package main

type APIGatewayAuthorizerRequest struct {
	Version               string              `json:"version"`
	Type                  string              `json:"type"`
	RouteArn              string              `json:"routeArn"`
	RouteKey              string              `json:"routeKey"`
	IdentitySource        []string            `json:"identitySource"`
	//RawQueryString        string              `json:"rawQueryString"`
	//RawPath               string              `json:"rawPath"`
	//Cookies               []string            `json:"cookies"`
	//Headers               map[string]string   `json:"headers"`
	//MultiValueHeaders     map[string][]string `json:"multiValueHeaders"`
	//QueryStringParameters map[string]string   `json:"queryStringParameters"`
	////MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"` // TODO RnD this field via jsonRawMessage
	PathParameters map[string]string                  `json:"pathParameters"`
	StageVariables map[string]string                  `json:"stageVariables"`
	RequestContext APIGatewayAuthorizerRequestContext `json:"requestContext"`
}

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