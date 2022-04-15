package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// region and session initialize together
func createSessionWithConfig() (*session.Session, error) {
	session_sqs, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
	})

	return session_sqs, err
}

// get SQS queue url for sending message
func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// this funcion is work for sending message to SQS
func SendSqsMessage(sess *session.Session, queueUrl string, messageBody string) error {
	sqsClient := sqs.New(sess)

	if flagValue != "" {
		_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(10),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"Flag": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(flagValue),
				},
			},
			QueueUrl:    &queueUrl,
			MessageBody: aws.String(messageBody),
		})
		return err
	} else {
		_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(10),
			QueueUrl:     &queueUrl,
			MessageBody:  aws.String(messageBody),
		})
		return err
	}

}

// this function will handle the SQS message processing and error controlling
func sqsHandler(mainConfig Configuration, alb_event events.ALBTargetGroupRequest) (MyResponse, error) {

	sess, _ := createSessionWithConfig()

	queueName := mainConfig.SqsName
	messageBody := alb_event.Body

	urlRes, err := GetQueueURL(sess, queueName)

	if err != nil {
		err = handleDeadLetterQueue(messageBody, mainConfig.DeadLetterQueueName)
		if err != nil {

			return makeResponse("Got an error while trying to get queue url as well as Dead letter queue", headers, 400), nil

		}
		fmt.Printf("Got an error while trying to get queue url: %v", err)
		return makeResponse("Got an error while trying to get queue url. However, send message to DLQ successfully", headers, 200), nil

	}

	err = SendSqsMessage(sess, *urlRes.QueueUrl, messageBody)

	if err != nil {
		log.Printf("Got an error while trying to send message to queue: %v", err)
		Logger("Got an error while trying to send message to queue", false, messageBody)

		err = handleDeadLetterQueue(messageBody, mainConfig.DeadLetterQueueName)

		if err != nil {
			return makeResponse("Got an error while trying to send message to queue as well as Dead letter queue", headers, 400), nil
		}

		return makeResponse("Got an error while trying to send message to queue. However, send message to DLQ successfully", headers, 200), nil

	}

	Logger("Message sent successfully", true, messageBody)

	return makeResponse("Message sent successfully", headers, 200), nil
}
