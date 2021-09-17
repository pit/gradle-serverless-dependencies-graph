variable "name_prefix" {
  type        = string
  description = "Lambda function names prefix"
  default     = "gradle-dependencies"
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

variable "create_prod_domain" {
  type    = bool
  default = true
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

variable "users" {
  type = map(map(string))
}

variable "tags" {
  type        = map(any)
  description = "Additional tags to add to resources"
  default     = {}
}
