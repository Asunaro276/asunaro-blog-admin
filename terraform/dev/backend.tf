terraform {
  backend "s3" {
    bucket = "tfstate-${var.environment}-nakano"
    key    = "${var.name}-${var.environment}.tfstate"
    region = "ap-northeast-1"
  }
}
