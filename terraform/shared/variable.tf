variable "name" {
  type        = string
  description = "The name of the project"

  validation {
    condition     = can(regex("^[a-zA-Z0-9-]+$", var.name))
    error_message = "The name must contain only letters, numbers and hyphens."
  }
}

variable "environment" {
  type        = string
  description = "The environment"

  validation {
    condition     = var.environment == "qa" || var.environment == "production"
    error_message = "The environment must be either 'qa' or 'production'."
  }
}

variable "domain" {
  type        = string
  description = "The domain of the project"

  validation {
    condition     = can(regex("^[a-zA-Z0-9-]*$", var.domain))
    error_message = "The domain must contain only letters, numbers and hyphens."
  }
}

variable "language" {
  type        = string
  description = "The language of the project"

  validation {
    condition     = var.language == "go"
    error_message = "The language must be 'go'."
  }
}

locals {
  codedir_local_path              = "${path.module}/../../cms_api/${var.language}"
  package_local_path              = "${local.codedir_local_path}/main.zip"
  package_base64sha256_local_path = "${local.codedir_local_path}/main.base64sha256"
  package_s3_key                  = "cms/main.zip"
  package_base64sha256_s3_key     = "${local.package_s3_key}.base64sha256.txt"
}
