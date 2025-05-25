resource "aws_lambda_function" "cms" {
  function_name    = "${var.name}-api-lambda"
  description      = "Lambda function for the CMS API"
  s3_bucket        = aws_s3_bucket.lambda_assets.bucket
  s3_key           = local.package_s3_key
  handler          = "main"
  role             = aws_iam_role.cms_lambda.arn
  runtime          = "provided.al2023"
  source_code_hash = data.aws_s3_object.lambda_package_hash.body
  timeout          = 30
  environment {
    variables = {
      LOG_LEVEL = local.cloudwatch_log_level # 環境別にログレベルを設定
    }
  }

  vpc_config {
    subnet_ids         = [for subnet in aws_subnet.private : subnet.id]
    security_group_ids = [aws_security_group.lambda.id]
  }
}

resource "aws_iam_role" "cms_lambda" {
  name = "lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy" "cms_lambda" {
  role = aws_iam_role.cms_lambda.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Effect   = "Allow"
        Resource = "${aws_cloudwatch_log_group.cms.arn}:*"
      },
      {
        Action = [
          "ec2:CreateNetworkInterface",
          "ec2:DescribeNetworkInterfaces",
          "ec2:DeleteNetworkInterface"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}

resource "aws_cloudwatch_log_group" "cms" {
  retention_in_days = local.cloudwatch_retention_days
  name              = "/aws/lambda/${aws_lambda_function.cms.function_name}"

  tags = {
    Name     = "${var.name}-lambda-logs"
    Purpose  = "cost-optimized"
    LogLevel = local.cloudwatch_log_level
  }
}

resource "aws_security_group" "lambda" {
  name        = "lambda-sg"
  description = "Security group for Lambda function"
  vpc_id      = aws_vpc.main.id


  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "lambda-sg"
  }
}

variable "lambda_subnet_ids" {
  description = "List of subnet IDs where the Lambda function will be deployed"
  type        = list(string)
}
