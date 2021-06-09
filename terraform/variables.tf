variable "name_prefix" {
  type        = string
  description = "Lambda function names prefix"
  default     = "terraform-registry"
}

variable "distrib_path" {
  type        = string
  description = "Path to folder with lambda binaries"
}

variable "domain_name" {
  type = string
}

variable "domain_acm_arn" {
  type = string
}

variable "create_dev_domain" {
  type    = bool
  default = false
}

variable "domain_dev_name" {
  type    = string
  default = null
}

variable "domain_dev_acm_arn" {
  type    = string
  default = null
}

variable "access_log_arns" {
  type = object({
    dev  = string
    prod = string
  })
}

variable "access_log_format" {
  type    = string
  default = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.integrationErrorMessage"
}


variable "bucket_name" {
  type = string
}

variable "bucket_readonly_arns" {
  type = list(string)
}

variable "bucket_readwrite_arns" {
  type = list(string)
}

variable "bucket_admin_arns" {
  type = list(string)
}

variable "tags" {
  type        = map(any)
  description = "Additional tags to add to resources"
  default     = {}
}
