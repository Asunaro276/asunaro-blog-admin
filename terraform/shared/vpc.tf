locals {
  private_subnets = {
    "a" = {
      cidr_block        = "10.0.11.0/24"
      availability_zone = "ap-northeast-1a"
    }
    "c" = {
      cidr_block        = "10.0.12.0/24"
      availability_zone = "ap-northeast-1c"
    }
  }
}

# VPC作成
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "${var.name}-vpc"
  }
}

# プライベートサブネット
resource "aws_subnet" "private" {
  for_each = local.private_subnets

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value.cidr_block
  availability_zone = each.value.availability_zone

  tags = {
    Name = "${var.name}-private-subnet-${each.key}"
  }
}

# ルートテーブル (プライベート)
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.name}-private-route-table"
  }
}

# プライベートサブネットとルートテーブルの関連付け
resource "aws_route_table_association" "private" {
  for_each = aws_subnet.private

  subnet_id      = each.value.id
  route_table_id = aws_route_table.private.id
}
