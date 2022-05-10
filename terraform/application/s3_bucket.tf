
# this will create a single bucket
resource "aws_s3_bucket" "bucket" {
  bucket = "${var.bucket_name}-${var.environment}"
  force_destroy = true
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.s3_event_for_metadata_add.arn
    events              = ["s3:ObjectCreated:Put","s3:ObjectCreated:Post","s3:ObjectCreated:CompleteMultipartUpload"]
  }

  depends_on = [
    aws_lambda_permission.allow_bucket1
  ]
}
