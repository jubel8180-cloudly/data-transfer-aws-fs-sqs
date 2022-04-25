package main

import (
	"context"
	"encoding/json"
	"errors"

	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Message struct {
	MessageId string `json:"MessageId"`
	Message   string `json:"Message"`
	Attribute string `json:"Attribute"`
}

type LogMessage struct {
	Message string
	Error   string
	Status  bool
	Event   string
}

func HandleRequest(ctx context.Context, event events.SQSEvent) (string, error) {

	if len(event.Records) <= 0 {
		log.Printf("EVENT: %s", "data not available")
		return "no data found", errors.New("error: data not available")
	}

	success := false
	// this is working for getting all queue records data form SQS.
	for _, message := range event.Records {
		msg := Message{MessageId: message.MessageId, Message: message.Body, Attribute: message.Attributes["flag"]}

		var logMessage LogMessage

		err := json.Unmarshal([]byte(msg.Message), &logMessage)

		if err != nil {
			return "msg unmarshal error from DLQ", err
		}



		log.Printf("Message: %s", logMessage.Message)
		log.Printf("Error: %v", logMessage.Error)
		log.Printf("Status: %v", logMessage.Status)

		eventJson, _ := json.MarshalIndent(logMessage.Event, "", "  ")
		log.Printf("EVENT: %s", eventJson)

		success = true

	}

	if success {
		return "log create successfully", nil
	}

	return "log create error", errors.New("log create error from DLQ")
}

func main() {
	lambda.Start(HandleRequest)
}
