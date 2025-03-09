terraform {
  backend "s3" {
    bucket = "tfstate-nakano"
    key    = "qa/cms.tfstate"
    region = "ap-northeast-1"
  }
}
