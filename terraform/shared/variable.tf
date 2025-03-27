variable "name" {
  type        = string
  description = "The name of the project"

  validation {
    condition     = can(regex("^[a-zA-Z0-9-]+$", var.name))
    error_message = "The name must contain only letters, numbers and hyphens."
  }
}

variable "env" {
  type        = string
  description = "The environment"

  validation {
    condition     = var.env == "dev" || var.env == "qa" || var.env == "production"
    error_message = "The environment must be 'dev', 'qa' or 'production'."
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

variable "image_name" {
  type        = string
  description = "The name of the image"

  validation {
    condition     = can(regex("^[a-zA-Z0-9-]+$", var.image_name))
    error_message = "The image name must contain only letters, numbers and hyphens."
  }
}

variable "skip_credentials_validation" {
  type        = bool
  description = "Skip credentials validation"
  default     = false
}

variable "skip_requesting_account_id" {
  type        = bool
  description = "Skip requesting account id"
  default     = false
}

variable "aws_endpoint_url_s3" {
  type        = string
  description = "The endpoint url of the s3"
  default     = ""
}

locals {
  codedir_local_path              = "${path.module}/../../cms_api/${var.language}"
  package_local_path              = "${local.codedir_local_path}/output/main.zip"
  package_base64sha256_local_path = "${local.codedir_local_path}/output/main.base64sha256"
  package_s3_key                  = "cms/main.zip"
  package_base64sha256_s3_key     = "${local.package_s3_key}.base64sha256.txt"
}
