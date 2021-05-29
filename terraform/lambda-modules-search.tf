module "lambda_modules_search" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.0.0"

  function_name = "${var.name_prefix}-modules-search"
  description   = "Registry API: :namespace/search"
  handler       = "modules-search"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = {
    BUCKET_NAME = aws_s3_bucket.this.id
  }

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  source_path = "${var.distrib_path}/modules-search"

  tags = merge({
    Name = "${var.name_prefix}-modules-search"
  }, var.tags)
}
