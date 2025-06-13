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
    origin_request_policy_id   = aws_cloudfront_origin_request_policy.allow_querystring.id
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

resource "aws_route53_zone" "viewer" {
  name = var.viewer_domain
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

resource "aws_route53_record" "viewer_cloudfront" {
  zone_id = aws_route53_zone.viewer.zone_id
  name    = var.viewer_domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.viewer.domain_name
    zone_id                = aws_cloudfront_distribution.viewer.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_acm_certificate" "viewer" {
  domain_name       = var.viewer_domain
  validation_method = "DNS"
  provider          = aws.n_virginia
}

resource "aws_acm_certificate_validation" "viewer" {
  certificate_arn         = aws_acm_certificate.viewer.arn
  validation_record_fqdns = [for record in aws_route53_record.viewer : record.fqdn]
  provider                = aws.n_virginia
}
