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

  attribute {
    name = "IsPrivate"
    type = "S"
  }

  global_secondary_index {
    name            = "IsPrivate-Timestamp-index"
    hash_key        = "IsPrivate"
    range_key       = "Timestamp"
    projection_type = "ALL"
  }
}
