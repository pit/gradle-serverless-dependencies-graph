module "lambda_modules_download" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.0.0"

  function_name = "${var.name_prefix}-modules-download"
  description   = "Registry API: /:namespace/:name/:provider/:version/download"
  handler       = "modules-download"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  source_path = "${var.distrib_path}/modules-download"

  tags = merge({
    Name = "${var.name_prefix}-modules-download"
  }, var.tags)
}
