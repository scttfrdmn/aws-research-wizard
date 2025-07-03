variable "domain_name" {
  description = "Name of the research domain"
  type        = string
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "environment" {
  description = "Environment name (localstack, dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "spack_packages" {
  description = "List of Spack packages to install"
  type        = list(string)
  default     = []
}

variable "monthly_budget" {
  description = "Monthly budget in USD"
  type        = number
  default     = 100
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
