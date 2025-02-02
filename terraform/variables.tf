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

variable "basic_id" {
  type        = string
  description = "Basic auth ID"
}

variable "basic_pw" {
  type        = string
  description = "Basic auth password"
  sensitive   = true
}
