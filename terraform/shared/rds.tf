resource "aws_rds_cluster" "aurora_mysql_serverless" {
  cluster_identifier      = "aurora-mysql-serverless"
  engine                  = "aurora-mysql"
  engine_mode             = "serverless"
  database_name           = "mydb"
  master_username         = "admin"
  master_password         = var.database_password
  backup_retention_period = 7
  preferred_backup_window = "02:00-03:00"
  skip_final_snapshot     = true

  scaling_configuration {
    auto_pause               = true
    max_capacity             = 16
    min_capacity             = 1
    seconds_until_auto_pause = 300
    timeout_action           = "ForceApplyCapacityChange"
  }

  vpc_security_group_ids = [aws_security_group.aurora_mysql.id]
  db_subnet_group_name   = aws_db_subnet_group.aurora_mysql.name
}

resource "aws_db_subnet_group" "aurora_mysql" {
  name       = "aurora-mysql-subnet-group"
  subnet_ids = var.database_subnet_ids

  tags = {
    Name = "Aurora MySQL subnet group"
  }
}

resource "aws_security_group" "aurora_mysql" {
  name        = "aurora-mysql-sg"
  description = "Allow MySQL traffic from Lambda only"
  vpc_id      = var.vpc_id

  ingress {
    description     = "MySQL from Lambda"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "aurora-mysql-sg"
  }
}

resource "aws_rds_cluster_parameter_group" "aurora_mysql" {
  name        = "aurora-mysql-parameter-group"
  family      = "aurora-mysql5.7"
  description = "Aurora MySQL parameter group"

  parameter {
    name  = "character_set_server"
    value = "utf8mb4"
  }

  parameter {
    name  = "character_set_client"
    value = "utf8mb4"
  }
}

# 必要な変数の定義
variable "database_password" {
  description = "Password for the master DB user"
  type        = string
  sensitive   = true
}

variable "vpc_id" {
  description = "VPC ID where the DB will be deployed"
  type        = string
}

variable "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  type        = string
}

variable "database_subnet_ids" {
  description = "List of subnet IDs where the DB can be deployed"
  type        = list(string)
}

# 出力値の定義
output "rds_cluster_endpoint" {
  description = "The cluster endpoint"
  value       = aws_rds_cluster.aurora_mysql_serverless.endpoint
}

output "rds_cluster_id" {
  description = "The ID of the cluster"
  value       = aws_rds_cluster.aurora_mysql_serverless.id
}
