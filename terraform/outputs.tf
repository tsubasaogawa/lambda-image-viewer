output "origin_bucket_arn" {
  value = module.origin.s3_bucket_arn
}

output "cloudfront_origin_domain_name" {
  value = aws_cloudfront_distribution.origin.domain_name
}

output "cloudfront_viewer_domain_name" {
  value = aws_cloudfront_distribution.viewer.domain_name
}

output "route53_origin_zone_id" {
  value = aws_route53_zone.origin.zone_id
}

output "route53_origin_cloudfront_fqdn" {
  value = aws_route53_record.origin_cloudfront.fqdn
}

output "route53_viewer_zone_id" {
  value = aws_route53_zone.viewer.zone_id
}

output "route53_viewer_cloudfront_fqdn" {
  value = aws_route53_record.viewer_cloudfront.fqdn
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.item.name
}
