module "lambda_providers_versions" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.0.0"

  function_name = "${var.name_prefix}-providers-versions"
  description   = "Registry API: :namespace/:type/versions"
  handler       = "providers-versions"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  source_path = "${var.distrib_path}/providers-versions"

  tags = merge({
    Name = "${var.name_prefix}-providers-versions"
  }, var.tags)
}
