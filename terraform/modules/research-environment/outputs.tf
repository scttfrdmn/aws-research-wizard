output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.research_node.id
}

output "public_ip" {
  description = "Public IP address of the instance"
  value       = aws_instance.research_node.public_ip
}

output "private_ip" {
  description = "Private IP address of the instance"
  value       = aws_instance.research_node.private_ip
}

output "security_group_id" {
  description = "ID of the security group"
  value       = aws_security_group.research.id
}

output "s3_bucket" {
  description = "S3 bucket for research data"
  value       = aws_s3_bucket.research_data.bucket
}

output "log_group" {
  description = "CloudWatch log group"
  value       = aws_cloudwatch_log_group.research_logs.name
}
