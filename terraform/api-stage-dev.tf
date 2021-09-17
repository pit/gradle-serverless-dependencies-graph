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
    Name = "gradle-dependencies/dev"
  })

  lifecycle {
    ignore_changes = [deployment_id]
  }
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
    Name = "gradle-dependencies/dev"
  }, var.tags)
}

resource "aws_apigatewayv2_api_mapping" "dev" {
  count = var.create_dev_domain ? 1 : 0

  api_id      = module.api.apigatewayv2_api_id
  domain_name = aws_apigatewayv2_domain_name.dev[0].id
  stage       = aws_apigatewayv2_stage.dev.id
}

