locals {
  api_routes = {
    # "GET /" = {
    #   lambda = module.lambda_index.lambda_function_name
    # },

    //    "GET /dependencies/{id}" = {
    //      lambda              = module.lambda_dependencies_list.lambda_function_name
    //      authorizer_required = false
    //    },
    //
    //    "GET /dependencies/{id}/{version}" = {
    //      lambda              = module.lambda_dependencies_list.lambda_function_name
    //      authorizer_required = false
    //    },
    //
    //    "GET /repos/{repo}" = {
    //      lambda              = module.lambda_repos_list.lambda_function_name
    //      authorizer_required = false
    //    },
    //
    //    "GET /repos/{repo}/{ref}" = {
    //      lambda              = module.lambda_repos_list.lambda_function_name
    //      authorizer_required = false
    //    },

    "PUT /repositories/v1/{org}/{repo}/{ref+}" = {
      lambda              = module.lambda_repo_batch_insert_put.lambda_function_name
      authorizer_required = false
    },

    "$default" = {
      lambda = module.lambda_default.lambda_function_name
    },
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

module "api" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "v1.0.0"

  name          = "gradle-dependencies"
  description   = "Serverless gradle dependencies app"
  protocol_type = "HTTP"

  # credentials_arn = module.apigateway_role.iam_role_arn
  # Custom domain
  # domain_name                      = var.domain_name
  # domain_name_certificate_arn      = var.domain_acm_arn
  create_api_domain_name           = false
  create_default_stage_api_mapping = true

  # Routes and integrations
  integrations = {
    for route_path in keys(local.api_routes) : (route_path) => merge(
      {
        lambda_arn             = "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:${lookup(local.api_routes[route_path], "lambda")}"
        credentials            = module.apigateway_role.iam_role_arn
        payload_format_version = "2.0"
        timeout_milliseconds   = 5000
      },
      # can(local.api_routes[route_path].authorizer_required) ? local.api_routes[route_path].authorizer_required ? {
      #   authorization_type = "CUSTOM"
      #   authorizer_id      = aws_apigatewayv2_authorizer.basic_auth.id
      # } : {} : {},
      can(local.api_routes[route_path].api_key_required) ? {
        api_key_required = local.api_routes[route_path].api_key_required
      } : {}
    )
  }

  tags = merge(var.tags, {
    Name = "terraform-registry"
  })
}
