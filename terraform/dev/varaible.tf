
variable "bucket_name" {
  type = string
  description = "Please provide the s3 bucket name"

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

variable "development_environment" {
  default = "dev"
}
variable "kinesis_processor_lambda_function_name" {
  default = "firehose_data_processor"
}
