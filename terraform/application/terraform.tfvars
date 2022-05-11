
region = "ap-south-1"

environment = "dev"

sqs_fs_lambda_function_name = "firehose_sqs_postback_event"

receive_dlq_lambda_function_name = "receive_event_message_from_dlq"

s3_event_function_name = "metadata_and_checksum_add_with_s3_event"

delivery_stream_name = "new-s3-lambda-delivery-stream"

sqs_name = "sqs_main"

dead_letter_queue_name = "dead_letter_queue_main"

load_balancer_name = "lb-firehose-sqs"

kinesis_processor_lambda_function_name = "firehose_data_processor_new"

ingress_cidr_blocks = ["37.111.207.253/32","27.147.205.239/32"]



