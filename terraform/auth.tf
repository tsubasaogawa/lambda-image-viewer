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

resource "random_password" "salt_for_private_image" {
  length  = 32
  special = true
}

resource "aws_secretsmanager_secret" "salt_for_private_image" {
  name        = "${var.origin_domain}-salt-for-private-image"
  description = "Token for private image access"
}

resource "aws_secretsmanager_secret_version" "salt_for_private_image" {
  secret_id     = aws_secretsmanager_secret.salt_for_private_image.id
  secret_string = random_password.salt_for_private_image.result
}
