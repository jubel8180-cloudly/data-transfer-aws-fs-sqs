# display necessary outputs
output "region"{
  value = var.region
}

output "lambda_function_names" {
  value = [aws_lambda_function.main.function_name,aws_lambda_function.receive_msg_from_dlq.function_name,aws_lambda_function.kinesis_lambda_processor.function_name]
}

output "s3_bucket_name"{
  value = "${aws_s3_bucket.bucket.id}"
}

output "delivery_stream" {
  value = aws_kinesis_firehose_delivery_stream.extended_s3_stream.name
}

output "sqs_dead_letter_queue_name" {
  value = aws_sqs_queue.terraform_queue_deadletter.name
}

output "sqs_name" {
  value = aws_sqs_queue.sqs_main.name
}

data "aws_caller_identity" "current" {}


output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "load_balancer_name" {
  value = module.my_ec2.load_balancer_name
}

output "load_balncer_dns_url"{
  value = module.my_ec2.aws_lb_dns_name
}


