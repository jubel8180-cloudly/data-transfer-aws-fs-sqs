package main_test

import (
	"context"
	"fmt"
	main "s3_event_notification"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	res, err := main.Handler(context.TODO(), events.S3Event{
		Records: []events.S3EventRecord{
			events.S3EventRecord{
				EventVersion:      "",
				EventSource:       "",
				AWSRegion:         *aws.String("ap-south-1"),
				EventTime:         time.Time{},
				EventName:         "",
				PrincipalID:       events.S3UserIdentity{},
				RequestParameters: events.S3RequestParameters{},
				ResponseElements:  map[string]string{},
				S3: events.S3Entity{
					SchemaVersion:   "1.0",
					ConfigurationID: "123123",
					Bucket: events.S3Bucket{
						Name: "se",
						OwnerIdentity: events.S3UserIdentity{
							PrincipalID: "12312312",
						},
						Arn: "123123123",
					},
					Object: events.S3Object{
						Key:       "asdasd",
						Size:      1213,
						Sequencer: "123123",
					},
				},
			},
		},
	})
	assert.Equal(t, res, "asdasdasdasd")
	fmt.Println(res)
	fmt.Println(err)
}
