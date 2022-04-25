# display necessary outputs
output "region"{
  value = var.region
}

output "lambda_function_names" {
  value = [var.sqs_fs_lambda_function_name,var.receive_dlq_lambda_function_name]
}

output "s3_bucket_name"{
  value = var.bucket_name
}

output "delivery_stream" {
  value = var.delivery_stream_name
}

output "sqs_dead_letter_queue_name" {
  value = var.dead_letter_queue_name
}

output "sqs_name" {
  value = var.sqs_name
}

data "aws_caller_identity" "current" {}


output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "load_balancer_name" {
  value = var.load_balancer_name
}

output "load_balncer_dns_url"{
  value = module.my_ec2.aws_lb_dns_name
}


