package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
		
)

// global variable
var region string
var flagValue string = ""

var headers Headers

// this is header setup for response
type Headers struct {
	ContentType string `json:"Content-Type"`
}

// this is the response representation
type MyResponse struct {
	Body            string  `json:"body"`
	IsBase64Encoded bool    `json:"isBase64Encoded"`
	StatusCode      int     `json:"statusCode"`
	Header          Headers `json:"headers"`
}

// we will take configuration from environment
type Configuration struct {
	DeliveryStreamName  string
	SqsName             string
	DeadLetterQueueName string
}

// This is simple response function to make a client back response.
func makeResponse(body string, headers Headers, status_code int) MyResponse {
	response := MyResponse{
		Body:            body,
		IsBase64Encoded: false,
		StatusCode:      status_code,
		Header:          headers,
	}
	return response
}

// this is simple logger function which will log to the cloudwatch
func Logger(status string, success bool, messageBody string) {
	log.Printf("STATUS: %s", status)
	log.Printf("MESSAGE_SUCCESS: %t", success)
	log.Printf("MESSAGE: %s", messageBody)
}

func isConditionFlagAvailable(jsonData  map[string]interface{},alb_event events.ALBTargetGroupRequest,mainConfig Configuration) (string, int, bool, map[string]interface{}){
	
	var conditionJson map[string]interface{} = nil


	if _, ok := jsonData["condition"]; !ok {

		Logger("Condition key not found", false, fmt.Sprintf("%s", jsonData))

		err := handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return "Condition key not found as well as Dead letter queue send error",400,true,conditionJson
			// return makeResponse("Condition key not found as well as Dead letter queue send error", headers, 400), nil

		}
		return "Condition key not found in payload. However, send message to Dead letter queue successfully delivered.",200,true,conditionJson
		// return makeResponse("Condition key not found in payload. However, send message to Dead letter queue successfully delivered.", headers, 200), nil
	}

	conditionJson = jsonData["condition"].(map[string]interface{})

	if _, ok := conditionJson["flag"]; !ok {

		Logger("Condition flag not found", false, fmt.Sprintf("%s", jsonData))

		err := handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return "Condition flag not found as well as Dead letter queue send error",400,true,conditionJson
		}
		return "Condition flag not found. However, send message to Dead letter queue successfully delivered.",200,true,conditionJson	}

	return "",200,false,conditionJson
}

// this is main handler function which will control SQS and firehose data tranfer based on condition
func Handler(ctx context.Context, alb_event events.ALBTargetGroupRequest) (MyResponse, error) {

	// providing correct response header format
	headers = Headers{ContentType: "application/json"}

	region = os.Getenv("region")
	// we are taking region and firehose delivery stream name from environment which is already setup in the lambda function
	mainConfig := Configuration{
		DeliveryStreamName:  os.Getenv("delivery_stream_name"),
		SqsName:             os.Getenv("main_sqs_name"),
		DeadLetterQueueName: os.Getenv("dead_letter_queue_name"),
	}

	// check request data is empty or not. if empty then return nil without putting data to kinesis
	if alb_event.Body == "" {

		Logger("Data not found", false, "")

		return makeResponse("Please provide a payload!", headers, 400), nil

	}

	// convert json data to byte format. Without byte format we can not pass data into kinesis firehose

	jsonData := make(map[string]interface{})

	err := json.Unmarshal([]byte(alb_event.Body), &jsonData)

	if err != nil {

		Logger("Payload Data is not valid", false, fmt.Sprintf("%s", jsonData))

		err := handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return makeResponse("Payload Data format error as well as Dead letter queue send error", headers, 400), nil

		}
		return makeResponse("Payload Data format error. However, Dead letter queue send successfully.", headers, 200), nil

	}

	if len(jsonData) == 0 {
		return makeResponse("Payload is empty!", headers, 400), nil
	}


	status_msg, status_code, new_err, conditionJson := isConditionFlagAvailable(jsonData,alb_event,mainConfig)

	if new_err{
		return makeResponse(status_msg, headers, status_code), nil
	}

	flagValue = strings.ToUpper(fmt.Sprintf("%v", conditionJson["flag"]))

	if flagValue == "A" {
		return firehoseHandler(mainConfig, alb_event)
	} else if flagValue == "Y" {
		return sqsHandler(mainConfig, alb_event)
	} else {

		Logger("Condition flag not valid", false, fmt.Sprintf("%s", jsonData))

		err := handleDeadLetterQueue(alb_event.Body, mainConfig.DeadLetterQueueName)
		if err != nil {
			return makeResponse("Data format error as well as Dead letter queue send error", headers, 400), nil

		}
		return makeResponse("Condition flag value "+flagValue+" is not valid. However, send message to Dead letter queue successfully delivered.", headers, 200), nil

	}

}


func main() {
	lambda.Start(Handler)
}

