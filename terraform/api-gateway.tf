locals {
  api_routes = {
    "GET /" = {
      lambda = module.lambda_index.lambda_function_name
    },

    "GET /.well-known/terraform.json" = {
      lambda              = module.lambda_discovery.lambda_function_name
      authorizer_required = true
    },

    # TODO Phase 3
    # "GET /modules/v1" = {
    #   lambda = module.lambda_modules_list.lambda_function_name
    # },
    # TODO Phase 3
    # "GET /modules/v1/{namespace}" = {
    #   lambda = module.lambda_modules_list.lambda_function_name
    # },

    # TODO Phase 3
    # "GET /modules/v1/search" = {
    #   lambda = module.lambda_modules_search.lambda_function_name
    # },

    "GET /modules/v1/{namespace}/{name}/{provider}/versions" = {
      lambda              = module.lambda_modules_versions.lambda_function_name
      authorizer_required = true
    },

    # Unknown URL
    # "GET /modules/v1/{namespace}/{name}/{provider}/download" = {
    #   lambda = module.lambda_modules_download.lambda_function_name
    # },

    "GET /modules/v1/{namespace}/{name}/{provider}/{version}/download" = {
      lambda              = module.lambda_modules_download.lambda_function_name
      authorizer_required = true
    },
    # TODO Phase 2
    # TODO Implement module archive/metainfo upload
    # "POST /modules/custom/{namespace}/{name}/{provider}/{version}/upload" = {
    #   lambda = module.lambda_modules_upload.lambda_function_name
    # },

    # TODO Phase 3
    # "GET /modules/v1/{namespace}/{name}" = {
    #   lambda = module.lambda_modules_latest_version.lambda_function_name
    # },
    # TODO Phase 3
    # "GET /modules/v1/{namespace}/{name}/{provider}" = {
    #   lambda = module.lambda_modules_latest_version.lambda_function_name
    # },

    # TODO Phase 3
    # "GET /modules/v1/{namespace}/{name}/{provider}/{version}" = {
    #   lambda = module.lambda_modules_get.lambda_function_name
    # },


    "GET /providers/v1/{namespace}/{type}/versions" = {
      lambda              = module.lambda_providers_versions.lambda_function_name
      authorizer_required = true
    },

    "GET /providers/v1/{namespace}/{type}/{version}/download/{os}/{arch}" = {
      lambda              = module.lambda_providers_download.lambda_function_name
      authorizer_required = true
    },
    # TODO Phase 2
    # TODO Implement module archive/metainfo upload
    # "POST /providers/custom/{namespace}/{type}/{version}/upload/{os}/{arch}" = {
    # lambda = module.lambda_providers_upload.lambda_function_name
    # },

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

  name          = "terraform-registry"
  description   = "Serverless terraform registry API"
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
      can(local.api_routes[route_path].authorizer_required) ? {
        authorization_type = "CUSTOM"
        authorizer_id      = aws_apigatewayv2_authorizer.basic_auth.id
      } : {},
      can(local.api_routes[route_path].api_key_required) ? {
        api_key_required = local.api_routes[route_path].api_key_required
      } : {}
    )
  }

  tags = merge(var.tags, {
    Name = "terraform-registry"
  })
}
