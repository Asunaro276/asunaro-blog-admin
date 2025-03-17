terraform {
  backend "s3" {
    bucket = "tfstate-nakano"
    key    = "dev/cms.tfstate"
    region = "ap-northeast-1"
  }
}
