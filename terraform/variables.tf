variable "origin_domain" {
  type        = string
  description = "Origin (Lambda) domain name"
}

variable "viewer_domain" {
  type        = string
  description = "Viewer (CloudFront) domain name"
}

variable "lambda_url" {
  type        = string
  description = "Lambda original endpoint url"
}

variable "lambda_function_name" {
  type        = string
  description = "Viewer Lambda function name"
}

variable "cloudfront_public_key" {
  type        = string
  description = "PEM-encoded CloudFront URL signing public key"
}

variable "maintenance_web_acl_arn" {
  type        = string
  description = "Optional CloudFront WAF web ACL for restricting access during migration"
  default     = null
}

variable "basic_id" {
  type        = string
  description = "Basic auth ID"
}

variable "basic_pw" {
  type        = string
  description = "Basic auth password"
  sensitive   = true
}
