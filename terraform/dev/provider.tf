provider "aws" {
  s3_use_path_style = true
  skip_requesting_account_id = true
  skip_credentials_validation = true
  skip_metadata_api_check = true

  endpoints {
    s3 = var.aws_endpoint_url_s3
  }
}
