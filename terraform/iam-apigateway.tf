module "apigateway_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"
  version = "v4.0.0"

  role_name = "${var.name_prefix}-api-gateway"

  trusted_role_services = [
    "apigateway.amazonaws.com",
  ]

  create_role       = true
  role_requires_mfa = false

  custom_role_policy_arns = [
    aws_iam_policy.apigateway_policy.arn,
  ]
  number_of_custom_role_policy_arns = 1
}

data "aws_iam_policy_document" "apigateway_policy" {
  statement {
    sid = "CloudWatchLogs"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["*"] # TODO use more strict cloudwatch logs arn
  }

  statement {
    sid = "InvokeLambda"
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = concat(
      [for route_path, route_obj in local.api_routes : "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:${route_obj.lambda}"],
      [module.lambda_authorizer.lambda_function_arn],
    )
  }
}

resource "aws_iam_policy" "apigateway_policy" {
  name   = "${var.name_prefix}-apigateway"
  policy = data.aws_iam_policy_document.apigateway_policy.json
}
