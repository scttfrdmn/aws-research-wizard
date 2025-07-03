variable "domain_name" {
  description = "Name of the research domain"
  type        = string
  default     = "digital_humanities"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "spack_packages" {
  description = "List of Spack packages to install"
  type        = list(string)
  default     = ["python@3.11.5", "py-nltk@3.8.1", "py-spacy@3.6.1", "py-pandas@2.0.3", "git@2.41.0"]
}

variable "monthly_budget" {
  description = "Monthly budget in USD"
  type        = number
  default     = 750
}
