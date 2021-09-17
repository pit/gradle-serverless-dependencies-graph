output "target_domain_zone_id" {
  value = var.create_prod_domain ? aws_apigatewayv2_domain_name.prod[0].domain_name_configuration[0].hosted_zone_id : ""
}

output "target_domain_name" {
  value = var.create_prod_domain ? aws_apigatewayv2_domain_name.prod[0].domain_name_configuration[0].target_domain_name : ""
}

output "target_dev_domain_zone_id" {
  value = var.create_dev_domain ? aws_apigatewayv2_domain_name.dev[0].domain_name_configuration[0].hosted_zone_id : ""
}

output "target_dev_domain_name" {
  value = var.create_dev_domain ? aws_apigatewayv2_domain_name.dev[0].domain_name_configuration[0].target_domain_name : ""
}
