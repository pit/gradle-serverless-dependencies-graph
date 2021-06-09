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

resource "aws_apigatewayv2_domain_name" "prod" {
  count = var.create_prod_domain ? 1 : 0

  domain_name = var.domain_name

  domain_name_configuration {
    certificate_arn = var.domain_acm_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = merge({
    Name = "terraform-registry/prod"
  }, var.tags)
}


resource "aws_apigatewayv2_api_mapping" "prod" {
  api_id      = module.api.apigatewayv2_api_id
  domain_name = aws_apigatewayv2_domain_name.prod[0].id
  stage       = aws_apigatewayv2_stage.prod.id
}
