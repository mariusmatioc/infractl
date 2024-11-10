package lambda

const Lambda = `
resource "aws_lambda_function" "lambda" {
  function_name = "${local.name}-{{.FunctionName}}"
  role          = aws_iam_role.lambda_role.arn
  handler       = "{{.Handler}}"
  runtime       = "{{.Runtime}}"

  filename      = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  ephemeral_storage {
	size = {{.EphemeralStorage}}
  }

{{if .EnvsString}}
  environment {
    variables = {
      {{.EnvsString}}
    }
  }
{{end}}

  vpc_config {
	subnet_ids         = [data.aws_subnet.private_subnet.id, data.aws_subnet.private_subnet_2.id]
	security_group_ids = [aws_security_group.lambda.id]
  }

  logging_config {
	log_group = aws_cloudwatch_log_group.lambda.name
	log_format = "Text"
  }

  layers = [
    {{.LayersString}}
  ]
}

resource "aws_security_group" "lambda" {
  name   = "${local.name}-{{.FunctionName}}"
  vpc_id = data.aws_vpc.vpc.id

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_iam_role" "lambda_role" {
  name = "${local.name}-{{.FunctionName}}"

  assume_role_policy = data.aws_iam_policy_document.assume_role.json
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
  ]
  inline_policy {
    name = "permissions"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action   = ["logs:*", "ec2:*", "s3:*"]
          Effect   = "Allow"
          Resource = "*"
        },
      ]
    })
  }
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

{{if .ScheduleExpression}}
resource "aws_cloudwatch_event_rule" "lambda_rule_1" {
  name        = "${local.name}-{{.FunctionName}}-1"
  description = "Trigger Lambda periodically"
  schedule_expression = "{{.ScheduleExpression}}"
}

resource "aws_cloudwatch_event_target" "lambda_target_1" {
  rule      = aws_cloudwatch_event_rule.lambda_rule_1.name
  target_id = "${local.name}-{{.FunctionName}}"
  arn       = aws_lambda_function.lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_1" {
    statement_id = "AllowExecutionFromCloudWatch-1"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.lambda.function_name
    principal = "events.amazonaws.com"
    source_arn = aws_cloudwatch_event_rule.lambda_rule_1.arn
}
{{end}}

{{if .S3ObjectCreated}}
data "aws_s3_bucket" "bucket" {
  bucket = "{{.S3ObjectCreated}}"
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.bucket.arn
  source_account = data.aws_caller_identity.current.account_id
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = data.aws_s3_bucket.bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.lambda.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}
{{end}}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = "{{.SourceFolder}}"
  output_path = "{{.BuildFolder}}/lambda.zip"
}

resource "aws_cloudwatch_log_group" "lambda" {
  name = "${local.name}-{{.FunctionName}}"
  retention_in_days = 14
}

output "{{.FunctionName}}" {
  value = aws_lambda_function.lambda.arn		
}

`
