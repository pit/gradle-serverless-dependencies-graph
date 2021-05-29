output "bucket_name" {
  value = aws_s3_bucket.this.id
}

output "bucket_arn" {
  value = aws_s3_bucket.this.arn
}

output "target_domain_zone_id" {
  value = module.api.apigatewayv2_domain_name_hosted_zone_id
}

output "target_domain_name" {
  value = module.api.apigatewayv2_domain_name_target_domain_name
}

output "stage_dev_url" {
  value = aws_apigatewayv2_stage.dev.invoke_url
}
