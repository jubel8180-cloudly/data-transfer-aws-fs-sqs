package main

import (
	"encoding/json"
	"fmt"
)

// this function will work for handling dead letter queue.
func handleDeadLetterQueue(logMessage LogMessage, dead_letter_queue_name string) error {

	sess, err := createSessionWithConfig()

	if err != nil {
		fmt.Printf("Session created error for SQS: %v", err)

		return err
	}

	urlRes, err := GetQueueURL(sess, dead_letter_queue_name)

	if err != nil {
		return err

	}

	messageBody, _ := json.Marshal(logMessage)
	err = SendSqsMessage(sess, *urlRes.QueueUrl, string(messageBody))

	if err != nil {
		fmt.Printf("Got an error while trying to send message to queue: %v", err)

		return err
	}

	return nil
}
