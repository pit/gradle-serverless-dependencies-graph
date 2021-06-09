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

  # Custom domain
  domain_name                      = var.domain_name
  domain_name_certificate_arn      = var.domain_acm_arn
  create_api_domain_name           = true
  create_default_stage_api_mapping = false

  # Routes and integrations
  integrations = {
    for route_path in keys(local.api_routes) : (route_path) => merge(
      {
        lambda_arn             = "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:${lookup(local.api_routes[route_path], "lambda")}"
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

resource "aws_apigatewayv2_stage" "dev" {
  api_id      = module.api.apigatewayv2_api_id
  name        = "dev"
  auto_deploy = true

  stage_variables = {
    LOG_LEVEL = "TRACE"
  }

  access_log_settings {
    destination_arn = var.access_log_arns.dev
    format          = var.access_log_format
  }

  default_route_settings {
    detailed_metrics_enabled = true
    throttling_burst_limit   = 5000
    throttling_rate_limit    = 10000
  }

  dynamic "route_settings" {
    for_each = { for key in keys(local.api_routes) : key => local.api_routes[key] if key != "$default" }
    content {
      route_key                = route_settings.key
      detailed_metrics_enabled = true
      throttling_burst_limit   = 5000
      throttling_rate_limit    = 10000
    }
  }

  tags = merge(var.tags, {
    Name = "terraform-registry/dev"
  })

  lifecycle {
    ignore_changes = [deployment_id]
  }
}

resource "aws_apigatewayv2_stage" "prod" {
  api_id      = module.api.apigatewayv2_api_id
  name        = "prod"
  auto_deploy = false

  stage_variables = {
    LOG_LEVEL = "INFO"
  }

  access_log_settings {
    destination_arn = var.access_log_arns.prod
    format          = var.access_log_format
  }

  default_route_settings {
    detailed_metrics_enabled = true
    throttling_burst_limit   = 5000
    throttling_rate_limit    = 10000
  }

  dynamic "route_settings" {
    for_each = { for key in keys(local.api_routes) : key => local.api_routes[key] if key != "$default" }
    content {
      route_key                = route_settings.key
      detailed_metrics_enabled = true
      throttling_burst_limit   = 5000
      throttling_rate_limit    = 10000
    }
  }

  tags = merge(var.tags, {
    Name = "terraform-registry/prod"
  })

  lifecycle {
    ignore_changes = [deployment_id]
  }
}

resource "aws_apigatewayv2_api_mapping" "prod" {
  api_id      = module.api.apigatewayv2_api_id
  domain_name = module.api.apigatewayv2_domain_name_id
  stage       = aws_apigatewayv2_stage.prod.id
}

resource "aws_apigatewayv2_domain_name" "dev" {
  count = var.create_dev_domain ? 1 : 0

  domain_name = var.domain_dev_name

  domain_name_configuration {
    certificate_arn = var.domain_dev_acm_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = merge({
    Name = "terraform-registry/dev"
  }, var.tags)
}

resource "aws_apigatewayv2_api_mapping" "dev" {
  count = var.create_dev_domain ? 1 : 0

  api_id      = module.api.apigatewayv2_api_id
  domain_name = aws_apigatewayv2_domain_name.dev[0].id
  stage       = aws_apigatewayv2_stage.dev.id
}

