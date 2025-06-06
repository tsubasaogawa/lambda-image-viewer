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
