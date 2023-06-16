module "origin" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.13.0"

  bucket_prefix = "${local.project_name}-origin"
  acl           = "private"
  website = {
    index_document = "index.html"
    error_document = "error.html"
  }
}
