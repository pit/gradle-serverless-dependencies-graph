module "lambda_repository_list_by_parent" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.0.0"

  function_name = "${var.name_prefix}-web-repositories-list-by-parent"
  description   = "Gradle: GET /repository/{org}/{repo}"
  handler       = "web-repositories-list-by-parent"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = {
    DYNAMODB_TABLE_STORAGE      = aws_dynamodb_table.storage.id
    DYNAMODB_TABLE_REPOSITORIES = aws_dynamodb_table.repositories.id
    DYNAMODB_TABLE_DEPENDENCIES = aws_dynamodb_table.dependencies.id
  }

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  source_path = "${var.distrib_path}/web-repositories-list-by-parent"

  tags = merge({
    Name = "${var.name_prefix}-web-repositories-list-by-parent"
  }, var.tags)
}
