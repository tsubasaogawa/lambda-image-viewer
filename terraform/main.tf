module "origin" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.13.0"

  bucket                   = var.origin_domain
  block_public_policy      = false
  block_public_acls        = false
  ignore_public_acls       = false
  restrict_public_buckets  = false
  control_object_ownership = true
  object_ownership         = "BucketOwnerEnforced"

  website = {
    index_document = "index.html"
    error_document = "error.html"
  }

  attach_policy = false
}

resource "aws_cloudfront_origin_access_control" "origin" {
  name                              = "${var.origin_domain}-oac"
  description                       = "OAC for S3 origin bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "no-override"
  signing_protocol                  = "sigv4"
}

resource "aws_s3_object" "metadata_fetcher" {
  bucket       = module.origin.s3_bucket_id
  key          = "assets/metadata_fetcher.js"
  content      = replace(file("assets/metadata_fetcher.js.tftpl"), "$${TF_API_ENDPOINT}", var.viewer_domain)
  source_hash  = filemd5("assets/metadata_fetcher.js.tftpl")
  content_type = "text/javascript"
}

resource "aws_s3_object" "generate_access_token" {
  bucket       = module.origin.s3_bucket_id
  key          = "assets/generate_access_token.js"
  content      = file("assets/generate_access_token.js")
  source_hash  = filemd5("assets/generate_access_token.js")
  content_type = "text/javascript"
}

resource "aws_acm_certificate" "origin" {
  domain_name       = var.origin_domain
  validation_method = "DNS"
  provider          = aws.n_virginia
}

resource "aws_route53_zone" "origin" {
  name = var.origin_domain
}

resource "aws_route53_record" "origin" {
  for_each = {
    for dvo in aws_acm_certificate.origin.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 300
  type            = each.value.type
  zone_id         = aws_route53_zone.origin.zone_id
}

resource "aws_acm_certificate_validation" "origin" {
  certificate_arn         = aws_acm_certificate.origin.arn
  validation_record_fqdns = [for record in aws_route53_record.origin : record.fqdn]
  provider                = aws.n_virginia
}

resource "aws_cloudfront_distribution" "origin" {
  aliases             = [var.origin_domain]
  enabled             = true
  is_ipv6_enabled     = true
  http_version        = "http2"
  price_class         = "PriceClass_All"
  retain_on_delete    = false
  wait_for_deployment = true

  default_cache_behavior {
    allowed_methods = [
      "GET",
      "HEAD",
    ]
    cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6" # CachingOptimized
    cached_methods = [
      "GET",
      "HEAD",
    ]
    origin_request_policy_id   = "b689b0a8-53d0-40ab-baf2-68738e2966ac" # AllViewerExceptHostHeader
    response_headers_policy_id = "eaab4381-ed33-4a86-88ca-d9558dc6cd63" # CORS-with-preflight-and-SecurityHeadersPolicy
    smooth_streaming           = false
    target_origin_id           = module.origin.s3_bucket_id
    viewer_protocol_policy     = "allow-all"
  }

  ordered_cache_behavior {
    path_pattern = "*/private/*" # TODO: to envvar
    allowed_methods = [
      "GET",
      "HEAD",
    ]
    cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6" # CachingOptimized
    cached_methods = [
      "GET",
      "HEAD",
    ]
    origin_request_policy_id   = aws_cloudfront_origin_request_policy.allow_querystring.id
    response_headers_policy_id = "eaab4381-ed33-4a86-88ca-d9558dc6cd63" # CORS-with-preflight-and-SecurityHeadersPolicy
    smooth_streaming           = false
    target_origin_id           = module.origin.s3_bucket_id
    viewer_protocol_policy     = "allow-all"
    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.access_token_auth.arn
    }
  }

  origin {
    domain_name              = module.origin.s3_bucket_bucket_domain_name
    origin_id                = module.origin.s3_bucket_id
    origin_access_control_id = aws_cloudfront_origin_access_control.origin.id
  }

  restrictions {
    geo_restriction {
      locations        = []
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn            = aws_acm_certificate.origin.arn
    cloudfront_default_certificate = false
    minimum_protocol_version       = "TLSv1.2_2021"
    ssl_support_method             = "sni-only"
  }
}

resource "aws_s3_bucket_policy" "origin" {
  bucket = module.origin.s3_bucket_id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid    = "AllowCloudFrontServicePrincipal"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        },
        Action   = "s3:GetObject"
        Resource = "${module.origin.s3_bucket_arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.origin.arn
          }
        }
      }
    ]
  })
}

resource "aws_route53_record" "origin_cloudfront" {
  zone_id = aws_route53_zone.origin.zone_id
  name    = var.origin_domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.origin.domain_name
    zone_id                = aws_cloudfront_distribution.origin.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_zone" "viewer" {
  name = var.viewer_domain
}

resource "aws_acm_certificate" "viewer" {
  domain_name       = var.viewer_domain
  validation_method = "DNS"
  provider          = aws.n_virginia
}

resource "aws_route53_record" "viewer" {
  for_each = {
    for dvo in aws_acm_certificate.viewer.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 300
  type            = each.value.type
  zone_id         = aws_route53_zone.viewer.zone_id
}

resource "aws_acm_certificate_validation" "viewer" {
  certificate_arn         = aws_acm_certificate.viewer.arn
  validation_record_fqdns = [for record in aws_route53_record.viewer : record.fqdn]
  provider                = aws.n_virginia
}

resource "aws_cloudfront_distribution" "viewer" {
  aliases             = [var.viewer_domain]
  enabled             = true
  is_ipv6_enabled     = true
  http_version        = "http2"
  price_class         = "PriceClass_All"
  retain_on_delete    = false
  wait_for_deployment = true

  default_cache_behavior {
    allowed_methods = [
      "GET",
      "HEAD",
    ]
    cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6" # CachingOptimized
    cached_methods = [
      "GET",
      "HEAD",
    ]
    compress                   = true
    origin_request_policy_id   = "b689b0a8-53d0-40ab-baf2-68738e2966ac" # AllViewerExceptHostHeader
    response_headers_policy_id = "eaab4381-ed33-4a86-88ca-d9558dc6cd63" # CORS-with-preflight-and-SecurityHeadersPolicy
    smooth_streaming           = false
    target_origin_id           = var.lambda_url
    viewer_protocol_policy     = "allow-all"
  }

  ordered_cache_behavior {
    path_pattern = "/cameraroll/*" # TODO: to envvar
    allowed_methods = [
      "GET",
      "HEAD",
    ]
    cache_policy_id = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad" # CachingDisabled
    cached_methods = [
      "GET",
      "HEAD",
    ]
    compress                   = true
    origin_request_policy_id   = "b689b0a8-53d0-40ab-baf2-68738e2966ac" # AllViewerExceptHostHeader
    response_headers_policy_id = "eaab4381-ed33-4a86-88ca-d9558dc6cd63" # CORS-with-preflight-and-SecurityHeadersPolicy
    smooth_streaming           = false
    target_origin_id           = var.lambda_url
    viewer_protocol_policy     = "allow-all"
    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.basic_auth.arn
    }
  }

  origin {
    connection_attempts = 3
    connection_timeout  = 10
    domain_name         = var.lambda_url
    origin_id           = var.lambda_url

    custom_origin_config {
      http_port                = 80
      https_port               = 443
      origin_keepalive_timeout = 5
      origin_protocol_policy   = "https-only"
      origin_read_timeout      = 30
      origin_ssl_protocols = [
        "TLSv1.2",
      ]
    }
  }

  restrictions {
    geo_restriction {
      locations        = []
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn            = aws_acm_certificate.viewer.arn
    cloudfront_default_certificate = false
    minimum_protocol_version       = "TLSv1.2_2021"
    ssl_support_method             = "sni-only"
  }
}

resource "aws_cloudfront_function" "basic_auth" {
  name    = "${replace(var.origin_domain, "/[^a-zA-Z0-9-_]/", "-")}-basic-auth"
  runtime = "cloudfront-js-2.0"
  publish = true
  code = templatefile(
    "${path.module}/assets/basic_auth.js.tftpl",
    {
      authString = base64encode("${var.basic_id}:${var.basic_pw}")
    }
  )
}

resource "aws_cloudfront_function" "access_token_auth" {
  name    = "${replace(var.origin_domain, "/[^a-zA-Z0-9-_]/", "-")}-token-auth"
  runtime = "cloudfront-js-2.0"
  publish = true
  code = templatefile(
    "${path.module}/assets/access_token_auth.js.tftpl",
    {
      salt = aws_secretsmanager_secret_version.salt_for_private_image.secret_string
    }
  )
}

resource "aws_cloudfront_origin_request_policy" "allow_querystring" {
  name = "${replace(var.origin_domain, "/[^a-zA-Z0-9-_]/", "-")}-allow-querystring-policy"

  query_strings_config {
    query_string_behavior = "all"
  }

  cookies_config {
    cookie_behavior = "none"
  }

  headers_config {
    header_behavior = "none"
  }
}

resource "aws_dynamodb_table" "item" {
  name         = "${var.origin_domain}-item"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "N"
  }

  global_secondary_index {
    hash_key        = "Timestamp"
    name            = "Timestamp"
    projection_type = "KEYS_ONLY"
  }
}

data "aws_secretsmanager_random_password" "salt_for_private_image" {
  password_length     = 32
  exclude_numbers     = false
  exclude_punctuation = true
  include_space       = false
}

resource "aws_secretsmanager_secret" "salt_for_private_image" {
  name        = "${var.origin_domain}-salt-for-private-image"
  description = "Token for private image access"
}

resource "aws_secretsmanager_secret_version" "salt_for_private_image" {
  secret_id     = aws_secretsmanager_secret.salt_for_private_image.id
  secret_string = data.aws_secretsmanager_random_password.salt_for_private_image.random_password
}
