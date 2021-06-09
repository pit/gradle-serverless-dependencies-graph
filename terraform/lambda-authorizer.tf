module "lambda_authorizer" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.0.0"

  function_name = "${var.name_prefix}-authorizer"
  description   = "Registry API Authorizer"
  handler       = "authorizer"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = {
    USER_TEST = "test"
  }

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  source_path = "${var.distrib_path}/authorizer"

  tags = merge({
    Name = "${var.name_prefix}-authorizer"
  }, var.tags)
}

resource "aws_lambda_permission" "lambda_authorizer_permissions" {
  function_name = module.lambda_authorizer.lambda_function_name

  statement_id = "AllowInvokeFromApiGateway"
  action       = "lambda:InvokeFunction"
  principal    = "apigateway.amazonaws.com"

  source_arn = "${module.api.apigatewayv2_api_execution_arn}/authorizers/${aws_apigatewayv2_authorizer.basic_auth.id}"
}
