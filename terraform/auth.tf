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

resource "aws_cloudfront_function" "private_path_guard" {
  name    = "${replace(var.origin_domain, "/[^a-zA-Z0-9-_]/", "-")}-private-path-guard"
  runtime = "cloudfront-js-2.0"
  publish = true
  code    = file("${path.module}/assets/private_path_guard.js")
}

resource "aws_cloudfront_function" "private_cache_headers" {
  name    = "${replace(var.origin_domain, "/[^a-zA-Z0-9-_]/", "-")}-private-cache-headers"
  runtime = "cloudfront-js-2.0"
  publish = true
  code    = file("${path.module}/assets/private_cache_headers.js")
}

resource "aws_cloudfront_public_key" "private_image" {
  name_prefix = "${substr(replace(var.origin_domain, "/[^a-zA-Z0-9-_]/", "-"), 0, 80)}-private-"
  encoded_key = var.cloudfront_public_key

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_cloudfront_key_group" "private_image" {
  name  = "${replace(var.origin_domain, "/[^a-zA-Z0-9-_]/", "-")}-private-image"
  items = [aws_cloudfront_public_key.private_image.id]
}
