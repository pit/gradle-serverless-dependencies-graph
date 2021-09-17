resource "aws_apigatewayv2_authorizer" "basic_auth" {
  name = "terraform-registry-basic-auth"

  api_id           = module.api.apigatewayv2_api_id
  authorizer_type  = "REQUEST"
  authorizer_uri   = module.lambda_authorizer.lambda_function_invoke_arn
  identity_sources = ["$request.header.Authorization"]

  enable_simple_responses = false

  authorizer_credentials_arn        = module.apigateway_role.iam_role_arn
  authorizer_result_ttl_in_seconds  = 10 # 60 * 5
  authorizer_payload_format_version = "2.0"
}
