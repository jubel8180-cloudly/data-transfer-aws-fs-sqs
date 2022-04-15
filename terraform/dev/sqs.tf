resource "aws_sqs_queue" "terraform_queue_deadletter" {
  name                      = "${var.dead_letter_queue_name}"
  delay_seconds             = 10
  max_message_size          = 2048
  message_retention_seconds = 240
  receive_wait_time_seconds = 10
  tags = {
    Environment = "production"
  }
}

resource "aws_sqs_queue" "sqs_main" {
  name                      = "${var.sqs_name}"
  delay_seconds             = 10
  max_message_size          = 2048
  message_retention_seconds = 240
  receive_wait_time_seconds = 10
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.terraform_queue_deadletter.arn
    maxReceiveCount     = 4
  })
  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = ["${aws_sqs_queue.terraform_queue_deadletter.arn}"]
  })

  tags = {
    Environment = "production"
  }
}



resource "aws_lambda_event_source_mapping" "dlq_event" {
  event_source_arn = aws_sqs_queue.terraform_queue_deadletter.arn
  function_name    = aws_lambda_function.receive_msg_from_dlq.arn
}