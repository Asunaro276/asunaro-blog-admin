resource "aws_rds_cluster" "aurora_mysql_serverless" {
  cluster_identifier      = "aurora-mysql-serverless"
  engine                  = "aurora-mysql"
  engine_mode             = "provisioned"
  engine_version          = "8.0.mysql_aurora.3.04.1"
  database_name           = "mydb"
  master_username         = "admin"
  master_password         = var.database_password
  backup_retention_period = 7
  preferred_backup_window = "02:00-03:00"
  skip_final_snapshot     = true

  serverlessv2_scaling_configuration {
    max_capacity = 16
    min_capacity = 0.5
  }

  vpc_security_group_ids = [aws_security_group.aurora_mysql.id]
  db_subnet_group_name   = aws_db_subnet_group.aurora_mysql.name

  tags = {
    Name = "aurora-mysql-serverless-v2"
  }
}

resource "aws_db_subnet_group" "aurora_mysql" {
  name       = "aurora-mysql-subnet-group"
  subnet_ids = [for subnet in aws_subnet.private : subnet.id]

  tags = {
    Name = "Aurora MySQL subnet group"
  }
}

resource "aws_security_group" "aurora_mysql" {
  name        = "aurora-mysql-sg"
  description = "Allow MySQL traffic from Lambda only"
  vpc_id      = aws_vpc.main.id

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
  family      = "aurora-mysql8.0"
  description = "Aurora MySQL parameter group for Serverless v2"

  parameter {
    name  = "character_set_server"
    value = "utf8mb4"
  }

  parameter {
    name  = "character_set_client"
    value = "utf8mb4"
  }

  tags = {
    Name = "aurora-mysql-parameter-group-v2"
  }
}

variable "database_password" {
  description = "Password for the master DB user"
  type        = string
  sensitive   = true
}

variable "enable_reader" {
  description = "Enable Aurora reader instance"
  type        = bool
  default     = false
}
