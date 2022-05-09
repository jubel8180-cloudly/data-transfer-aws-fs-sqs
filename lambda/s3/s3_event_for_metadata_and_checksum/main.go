package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws/session"
)

func handler(ctx context.Context, s3Event events.S3Event) {

	var checkSumAlgo = "sha256"

	svc := s3.New(session.New())

	for _, record := range s3Event.Records {
		s3_obj := record.S3
		bucket_name := fmt.Sprintf("%v", s3_obj.Bucket.Name)
		key := fmt.Sprintf("%v", s3_obj.Object.Key)
		headInput := &s3.HeadObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(key),
		}

		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(key),
		})
		if err != nil {
			fmt.Println(err)
		}

		body, err := ioutil.ReadAll(obj.Body)
		if err != nil {
			fmt.Println(err)
		}

		reader := csv.NewReader(bytes.NewBuffer(body))
		record_data, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error", err)
		}

		total_length := len(record_data)

		for value := range record_data { // for i:=0; i<len(record)
			fmt.Println("", record_data[value])
		}

		sourceObject, _ := svc.HeadObject(headInput)
		// Error handling intentionally omitted
		meta := make(map[string]*string)
		for k, v := range sourceObject.Metadata {
			meta[k] = v
		}

		length_without_header := total_length - 1
		meta_data := fmt.Sprintf("%d", length_without_header)
		meta["total_number"] = &meta_data

		fmt.Println(meta)

		copyInput := &s3.CopyObjectInput{
			Bucket:            aws.String(bucket_name),
			CopySource:        aws.String(fmt.Sprintf("%v/%v", s3_obj.Bucket.Name, s3_obj.Object.Key)),
			Key:               aws.String(key),
			MetadataDirective: aws.String("REPLACE"),
			Metadata:          meta,
			ChecksumAlgorithm: &checkSumAlgo,
		}

		resp, err := svc.CopyObject(copyInput)

		if err != nil {
			fmt.Println("SUCCESS: TRUE")
			log.Println("Add metadata and checksum successfully")
		} else {
			fmt.Println("SUCCESS: FALSE")
			log.Println("metadata added failed")
		}

		fmt.Println(resp)

	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
