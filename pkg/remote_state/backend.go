package remote_state

const Backend = `
terraform {
  required_version = ">= 1.5.5"
}

provider "aws" {
}

resource "aws_s3_bucket" "remote_state" {
  bucket = "terraform-state-{{.ProjectName}}"
}

resource "aws_s3_bucket_versioning" "remote_state" {
  bucket = aws_s3_bucket.remote_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "remote_state" {
  bucket = aws_s3_bucket.remote_state.bucket
  rule {
	apply_server_side_encryption_by_default {
	  sse_algorithm = "AES256"
	}
  }
}

resource "aws_s3_bucket_public_access_block" "remote_state" {
  bucket = aws_s3_bucket.remote_state.bucket
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_dynamodb_table" "remote_state" {
  name           = "terraform-state-{{.ProjectName}}"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"
  attribute {
	name = "LockID"
	type = "S"
  }
}

`
