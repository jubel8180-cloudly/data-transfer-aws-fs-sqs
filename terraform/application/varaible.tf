
variable "bucket_name" {
  type = string
  description = "Please provide the s3 bucket name"
}

variable "environment" {
  type    = string
  default = "dev"
}

variable "region"{
  type = string
  description = "Please provide region name"
}

variable "sqs_fs_lambda_function_name" {
  default = "send_data_to_firehose_or_sqs"
}

variable "receive_dlq_lambda_function_name" {
  default = "receive_message_from_dlq"
}

variable "s3_event_function_name" {
  default = "metadata_and_checksum_add_with_s3_event"
}

variable "delivery_stream_name" {
  default = "new-s3-lambda-delivery-stream"
}

variable "sqs_name"{
  default = "dev_sqs_main"
}

variable "dead_letter_queue_name"{
  default = "dev_dead_letter_queue_main"
}
variable "load_balancer_name" {
  default = "lambda-firehose-s3-workflow-lb"
}

variable "kinesis_processor_lambda_function_name" {
  default = "firehose_data_processor"
}

variable "ingress_cidr_blocks" {
  description = "List of IPv4 CIDR ranges to use on all ingress rules"
  type        = list(string)
  default     = []
}
