terraform {
  backend "s3" {
    bucket = "tfstate-${var.env}-nakano"
    key    = "${var.name}-${var.env}.tfstate"
    region = "ap-northeast-1"
  }
}
