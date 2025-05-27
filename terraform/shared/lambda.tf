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
          "s3:GetObject"
        ]
        Effect   = "Allow"
        Resource = "${aws_s3_bucket.lambda_assets.arn}/*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_vpc_access" {
  role       = aws_iam_role.cms_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.cms_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_cloudwatch_log_group" "cms" {
  retention_in_days = 1
  name              = "/aws/lambda/${aws_lambda_function.cms.function_name}"

  tags = {
    Name     = "${var.name}-lambda-logs"
    Purpose  = "cost-optimized"
    LogLevel = "INFO"
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
