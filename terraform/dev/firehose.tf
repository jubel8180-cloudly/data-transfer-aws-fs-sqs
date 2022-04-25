
# create firehose stream
resource "aws_kinesis_firehose_delivery_stream" "extended_s3_stream" {
  name        = var.delivery_stream_name
  destination = "extended_s3"

  extended_s3_configuration {
    #  role_arn   = "arn:aws:iam::115391213665:role/lambda-execution-role-s3-access"
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.bucket.arn
    buffer_size        = 64
    buffer_interval    = 60
    # prefix = "kinesis_firehose_data/!{timestamp:yyyy}-!{timestamp:MM}-!{timestamp:dd}-!{timestamp:HH}/"
    prefix = "kinesis_firehose_data/!{partitionKeyFromQuery:app_id}/"
    error_output_prefix = "kinesis_firehose_error_data/!{firehose:random-string}/!{firehose:error-output-type}/!{timestamp:yyyy/MM/dd}/"


    dynamic_partitioning_configuration{
      enabled = "true"
    }

    processing_configuration {
      enabled = "true"

      # Multi-record deaggregation processor example
      processors {
        type = "RecordDeAggregation"
        parameters {
          parameter_name  = "SubRecordType"
          parameter_value = "JSON"
        }
      }

       processors {
        type = "MetadataExtraction"
        parameters {
          parameter_name  = "JsonParsingEngine"
          parameter_value = "JQ-1.6"
        }
        parameters {
          parameter_name  = "MetadataExtractionQuery"
          parameter_value = "{app_id:.app_id}"
        }
      }


    }
  }

}
