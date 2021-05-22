locals {
  lambdas_permissions = distinct([for route in values(local.api_routes) : route.lambda])
}
resource "aws_lambda_permission" "lambdas_permissions" {
  for_each = zipmap(local.lambdas_permissions, range(length(local.lambdas_permissions)))

  function_name = each.key

  statement_id = "AllowInvokeFromApiGateway"
  action       = "lambda:InvokeFunction"
  principal    = "apigateway.amazonaws.com"

  source_arn = "${module.api.apigatewayv2_api_execution_arn}/*/*/*"
}
