package main_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	main "github.com/jubel-cloudly/sqs-firehose-send-data"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		context context.Context
		request events.ALBTargetGroupRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			context: context.TODO(),
			request: events.ALBTargetGroupRequest{Body: "Paul"},
			expect:  "{\"msg\":\"Payload Data format error as well as Dead letter queue send error\",\"success\":false}",
			err:     nil,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			context: context.TODO(),
			request: events.ALBTargetGroupRequest{Body: ""},
			expect:  "{\"msg\":\"Please provide a payload!\",\"success\":true}",
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.context, test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}
