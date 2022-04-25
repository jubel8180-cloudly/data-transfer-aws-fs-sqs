package main

import (
	"encoding/json"
	"fmt"
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

// FakeEntity is just used for testing purposes
type FakeEntity struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func firehosePutRecordBatch(svc *firehose.Firehose, mainConfig Configuration, eventData []JsonEvent) (*firehose.PutRecordBatchOutput, error) {

	recordsBatchInput := &firehose.PutRecordBatchInput{}
	recordsBatchInput = recordsBatchInput.SetDeliveryStreamName(mainConfig.DeliveryStreamName)

	records := []*firehose.Record{}

	for _, data := range eventData {
		b, err := json.Marshal(data)

		if err != nil {
			log.Printf("Error: %v", err)
		}

		record := &firehose.Record{Data: b}
		records = append(records, record)
	}

	recordsBatchInput = recordsBatchInput.SetRecords(records)

	resp, err := svc.PutRecordBatch(recordsBatchInput)
	if err != nil {
		log.Printf("PutRecordBatch err: %v\n", err)
	} else {
		log.Printf("PutRecordBatch: %v\n", resp)
	}

	return resp, nil

}

// Records - convert [][]byte to kinesis types.Record
func Records(data [][]byte) []*firehose.Record {
	ret := make([]*firehose.Record, 0, len(data))
	for _, record := range data {
		ret = append(ret, &firehose.Record{Data: record})
	}
	return ret
}

// this function will work for firehose data tranfer controlling and error handling
func firehoseHandler(mainConfig Configuration, alb_event events.ALBTargetGroupRequest) (MyResponse, error) {
	sess := createSession()

	config := createConfig(mainConfig)

	var arr map[string]interface{}

	var eventData JsonEvents

	err := json.Unmarshal([]byte(alb_event.Body), &arr)

	if err != nil {

		log.Printf("event body unmarshal error: %v", err)
	}

	jsonString, _ := json.Marshal(arr["records"])

	err = json.Unmarshal([]byte(jsonString), &eventData)

	if err != nil {

		log.Printf("event body unmarshal error: %v", err)
	}

	if err != nil {

		log.Printf("Got an error while trying to send message to queue: %v", err)

		Logger("Json data format errors", false, string(jsonString))

		err = handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return makeResponse("Json Data format error as well as Dead letter queue send error", headers, 400), nil

		}

		return makeResponse("Json Data format error. However, send message to Dead letter queue successfully delivered.", headers, 400), nil
	}

	svc := firehose.New(sess, config)

	// put data into firehose
	// result, err := firehosePutRecord(svc, mainConfig, jsonString)
	result, err := firehosePutRecordBatch(svc, mainConfig, eventData)

	// check put record success or not
	if err != nil {
		log.Printf("Got an error while trying to send message to queue: %v", err)

		Logger("Json data format errors", false, string(jsonString))

		err = handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return makeResponse("Data format error as well as Dead letter queue send error", headers, 400), nil

		}
		return makeResponse("Data send error to firehose. However, send message to Dead letter queue successfully delivered.", headers, 400), nil
	}

	log.Println("Firehose Data Transfer and File created successfully")

	Logger("Firehose Data Transfer and File created successfully", true, string(jsonString))

	fmt.Println(result.RequestResponses)

	// make success response
	response := makeResponse("Firehose data transfer and file created successfully done.", headers, 200)

	return response, err
}
