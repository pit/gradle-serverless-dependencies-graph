module "lambda_repo_batch_insert_put" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.0.0"

  function_name = "${var.name_prefix}-repo-batch-insert-put"
  description   = "Gradle Dependencies: /:repo/:ref+"
  handler       = "repo-batch-insert-put"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = {
    DYNAMODB_TABLE = aws_dynamodb_table.gradle_dependencies.id
  }

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  source_path = "${var.distrib_path}/repo-batch-insert-put"

  tags = merge({
    Name = "${var.name_prefix}-repo-batch-insert-put"
  }, var.tags)
}
