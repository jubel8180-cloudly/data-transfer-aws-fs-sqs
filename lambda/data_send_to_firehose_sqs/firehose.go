package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

//  if we want to work with aws. Then we need to initializes a session using this function
func createSession() *session.Session {
	sess := session.Must(session.NewSession())

	return sess

}

// we need to setup in which region we will work on
func createConfig(mainConfig Configuration) *aws.Config {
	config := aws.NewConfig()

	config.WithRegion(region)

	return config
}

// passign data through firehose and this data will be stored in s3 bucket.
func firehosePutRecord(svc *firehose.Firehose, mainConfig Configuration, data []byte) (*firehose.PutRecordOutput, error) {

	result, err := svc.PutRecord(&firehose.PutRecordInput{
		DeliveryStreamName: aws.String(mainConfig.DeliveryStreamName),
		Record: &firehose.Record{
			Data: data,
		},
	})

	return result, err

}

// this function will work for firehose data tranfer controlling and error handling
func firehoseHandler(mainConfig Configuration, alb_event events.ALBTargetGroupRequest) (MyResponse, error) {
	sess := createSession()

	config := createConfig(mainConfig)

	jsonBodyData, err := json.Marshal(alb_event.Body)

	if err != nil {

		log.Printf("Got an error while trying to send message to queue: %v", err)

		Logger("Json data format errors", false, string(jsonBodyData))

		err = handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return makeResponse("Json Data format error as well as Dead letter queue send error", headers, 400), nil

		}

		return makeResponse("Json Data format error. However, send message to Dead letter queue successfully delivered.", headers, 400), nil
	}

	svc := firehose.New(sess, config)

	// put data into firehose
	result, err := firehosePutRecord(svc, mainConfig, jsonBodyData)

	// check put record success or not
	if err != nil {
		log.Printf("Got an error while trying to send message to queue: %v", err)

		Logger("Json data format errors", false, string(jsonBodyData))

		err = handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return makeResponse("Data format error as well as Dead letter queue send error", headers, 400), nil

		}
		return makeResponse("Data send error to firehose. However, send message to Dead letter queue successfully delivered.", headers, 400), nil
	}

	log.Println("Firehose Data Transfer and File created successfully")

	Logger("Firehose Data Transfer and File created successfully", true, string(jsonBodyData))

	log.Printf("RECORD_ID: %s", *result.RecordId)

	// make success response
	response := makeResponse("Firehose data transfer and file created successfully done.", headers, 200)

	return response, err
}
