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
    ]
  })
}

resource "aws_cloudwatch_log_group" "cms" {
  retention_in_days = 1
  name              = "/aws/lambda/${aws_lambda_function.cms.function_name}"
}
