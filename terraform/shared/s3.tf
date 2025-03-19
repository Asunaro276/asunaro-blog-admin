resource "aws_s3_bucket" "lambda_assets" {
  bucket = "lambda-assets-${var.name}-${var.env}"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "lambda_assets" {
  bucket = aws_s3_bucket.lambda_assets.bucket

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "lambda_assets" {
  bucket = aws_s3_bucket.lambda_assets.bucket

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_public_access_block" "lambda_assets" {
  bucket = aws_s3_bucket.lambda_assets.bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "null_resource" "lambda_package" {
  depends_on = [aws_s3_bucket.lambda_assets]

  provisioner "local-exec" {
    // Build the lambda package and zip it
    working_dir = "${local.codedir_local_path}/admin_api_server"
    environment = {
      GOOS   = "linux"
      GOARCH = "amd64"
    }
    command = "go build -o ./output/bootstrap ./cmd/main.go && cd output && zip main.zip bootstrap"
  }

  provisioner "local-exec" {
    // Upload the lambda package to the s3 bucket
    command = "aws s3 cp ${local.package_local_path} s3://${aws_s3_bucket.lambda_assets.bucket}/${local.package_s3_key}"
    environment = {
      AWS_ENDPOINT_URL_S3 = var.aws_endpoint_url_s3
    }
  }

  provisioner "local-exec" {
    // Generate the sha256 hash of the lambda package
    command = "openssl dgst -sha256 -binary ${local.package_local_path} | openssl enc -base64 | tr -d \"\n\" > ${local.package_base64sha256_local_path}"
  }

  provisioner "local-exec" {
    // Upload the sha256 hash of the lambda package to the s3 bucket
    command = "aws s3 cp ${local.package_base64sha256_local_path} s3://${aws_s3_bucket.lambda_assets.bucket}/${local.package_base64sha256_s3_key} --content-type \"text/plain\""
    environment = {
      AWS_ENDPOINT_URL_S3 = var.aws_endpoint_url_s3
    }
  }

  triggers = {
    code_diff = join("", [
      for file in fileset(local.codedir_local_path, "**/{*.go,*.mod,*.sum}") : filebase64sha256("${local.codedir_local_path}/${file}")
    ])
  }
}

data "aws_s3_object" "lambda_package" {
  depends_on = [null_resource.lambda_package, null_resource.bucket_empty]
  bucket     = aws_s3_bucket.lambda_assets.bucket
  key        = local.package_s3_key
}

data "aws_s3_object" "lambda_package_hash" {
  depends_on = [null_resource.lambda_package, null_resource.bucket_empty]
  bucket     = aws_s3_bucket.lambda_assets.bucket
  key        = local.package_base64sha256_s3_key
}

resource "null_resource" "bucket_empty" {
  triggers = {
    bucket = aws_s3_bucket.lambda_assets.bucket
  }
  provisioner "local-exec" {
    when    = destroy
    command = "aws s3 rm s3://${self.triggers.bucket} --recursive"
  }
}
