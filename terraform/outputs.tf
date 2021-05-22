output "target_domain_zone_id" {
  value = module.api.apigatewayv2_domain_name_hosted_zone_id
}

output "target_domain_name" {
  value = module.api.apigatewayv2_domain_name_target_domain_name
}
