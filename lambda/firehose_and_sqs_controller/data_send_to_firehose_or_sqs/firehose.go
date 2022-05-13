package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrorDeadLetterQueue = errors.New("data send error to dead letter queue")
)

//  if we want to work with aws. Then we need to initializes a session using this function
func createSession() *session.Session {
	sess := session.Must(session.NewSession())

	return sess

}

// we need to setup in which region we will work on
func createConfig() *aws.Config {
	config := aws.NewConfig()

	config.WithRegion(region)

	return config
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

	return resp, err

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

	config := createConfig()

	var arr map[string]interface{}

	var eventData JsonEvents

	err := json.Unmarshal([]byte(alb_event.Body), &arr)

	if err != nil {
		logMessage := LogMessage{
			Message: "Event body unmarshal error",
			Error:   fmt.Sprintf("%v", err),
			Status:  false,
			Event:   alb_event.Body,
		}
		err := handleDeadLetterQueue(logMessage, mainConfig.DeadLetterQueueName)
		if err != nil {
			return makeResponse("Event Data validation error as well as Dead letter queue send error", headers, 400), nil

		}

		log.Printf("event body unmarshal error: %v", err)

		return makeResponse("Event Data format error. However, send message to Dead letter queue successfully delivered.", headers, 200), nil

	}

	jsonString, _ := json.Marshal(arr["records"])

	err = json.Unmarshal([]byte(jsonString), &eventData)

	if err != nil {
		logMessage := LogMessage{
			Message: "Data records is not valid",
			Error:   fmt.Sprintf("%v", err),
			Status:  false,
			Event:   alb_event.Body,
		}
		err := handleDeadLetterQueue(logMessage, mainConfig.DeadLetterQueueName)
		log.Printf("event body unmarshal error: %v", err)

		if err != nil {
			return makeResponse("Event Data format error as well as Dead letter queue send error", headers, 400), nil

		}

		return makeResponse("Event Data format error. However, send message to Dead letter queue successfully delivered.", headers, 200), nil

	}

	svc := firehose.New(sess, config)

	// put data into firehose
	// result, err := firehosePutRecord(svc, mainConfig, jsonString)
	result, err := firehosePutRecordBatch(svc, mainConfig, eventData)

	// check put record success or not
	if err != nil {
		log.Printf("Got an error while trying to send message to queue: %v", err)

		Logger("Got an error while trying to send message to queue", false, string(jsonString))

		logMessage := LogMessage{
			Message: "Got an error while trying to send message to queue",
			Error:   fmt.Sprintf("%v", err),
			Status:  false,
			Event:   alb_event.Body,
		}
		err := handleDeadLetterQueue(logMessage, mainConfig.DeadLetterQueueName)

		if err != nil {
			return makeResponse("Data send error to firehose as well as Dead letter queue send error", headers, 400), ErrorDeadLetterQueue
		}

		return makeResponse("Data send error to firehose. However, send message to Dead letter queue successfully delivered.", headers, 200), nil
	} else {
		Logger("Firehose Data Transfer and File created successfully", true, string(jsonString))

		fmt.Println(result.RequestResponses)

		// make success response
		response := makeResponse("Firehose data transfer and file created successfully done.", headers, 200)

		return response, err
	}

}
