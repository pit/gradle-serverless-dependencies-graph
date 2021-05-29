module "lambdas_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"
  version = "v4.0.0"

  role_name = "${var.name_prefix}-lambdas"

  trusted_role_services = [
    "lambda.amazonaws.com",
  ]

  create_role       = true
  role_requires_mfa = false


  custom_role_policy_arns = [
    aws_iam_policy.lambdas_policy.arn,
  ]
  number_of_custom_role_policy_arns = 1
}

data "aws_iam_policy_document" "lambdas_policy" {
  statement {
    sid = "AllowReadObjects"
    actions = [
      "s3:GetObject",
      "s3:GetObjectAcl",
    ]
    resources = ["${aws_s3_bucket.this.arn}/*"]
  }

  statement {
    sid       = "AllowListObjects"
    actions   = ["s3:ListBucket"]
    resources = [aws_s3_bucket.this.arn]
  }

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
    resources = ["*"]
  }
}

resource "aws_iam_policy" "lambdas_policy" {
  name   = "${var.name_prefix}-lambdas"
  policy = data.aws_iam_policy_document.lambdas_policy.json
}

