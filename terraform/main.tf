module "origin" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.13.0"

  bucket                   = local.project_name
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

  attach_policy = true
  policy = jsonencode({
    version = "2012-10-17"
    statement = [{
      sid       = "PublicReadGetObject"
      effect    = "Allow"
      principal = "*"
      action    = "s3:GetObject"
      resource  = "${module.origin.s3_bucket_arn}/*"
    }]
  })
}
