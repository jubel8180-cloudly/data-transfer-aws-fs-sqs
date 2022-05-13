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
			request: events.ALBTargetGroupRequest{Body: "Hello"},
			expect:  "{\"msg\":\"Payload Data format error as well as Dead letter queue send error\",\"success\":false}",
			err:     nil,
		}, {
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			context: context.TODO(),
			request: events.ALBTargetGroupRequest{Body: `{"message":"Hello"}`},
			expect:  "{\"msg\":\"Condition key not found as well as Dead letter queue send error\",\"success\":false}",
			err:     nil,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			context: context.TODO(),
			request: events.ALBTargetGroupRequest{Body: ""},
			expect:  "{\"msg\":\"Please provide a payload!\",\"success\":true}",
			err:     nil,
		}, {
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			context: context.TODO(),
			request: events.ALBTargetGroupRequest{Body: `{
				"condition": {
				 "flag": "Y"
				 },
				 "records": {
					 "body": "All provider data collected successfully",
					 "attributes": {
						 "status": "success",
						 "provider": "dynata"

					 }
				 }
			 }`},
			expect: "{\"msg\":\"Got an error while trying to get queue url as well as Dead letter queue\",\"success\":false}",
			err:    nil,
		},
		{
			context: context.TODO(),
			request: events.ALBTargetGroupRequest{Body: `
			{
				"condition": {
				  "flag": "A"
				},
				"records": [
				  {
					"app_id":"app_1",
					"device_id":"608d99c2-9ff2-40a1-9ab5-8ee97861ccb2",
					"request_id":"19cc1f31-2a2e-4897-a8fd-a8dddf7a8e56"
				  },
				  {
					"app_id":"app_2",
					"device_id":"7d4586f0-e73f-49f0-ba9d-1b66f5413580",
					"request_id":"21805093-e285-40ed-b2c8-d0c7164f7953    "
				  }
				  ,{
					"app_id":"app_2",
					"device_id":"11fd3383-11f8-47f9-a110-37e2e409fdc9",
					"request_id":"f86b8e64-bc5f-4486-96ba-5792dc1836fd"
				  },
				  {
					"app_id":"app_1",
					"device_id":"887e6f18-2850-44df-a896-183403f94a63",
					"request_id":"d24c7916-fdf9-4275-8eac-926c16188615"
				  },
				  {
					"app_id":"app_1",
					"device_id":"07e140fa-ebb1-485f-848e-ab8d743f5788",
					"request_id":"4ed50b8b-5b47-4a62-a63a-2df073a3e602"
				  },
				  {
					"app_id":"app_2",
					"device_id":"83962919-5fbf-4790-8f17-ceeef6e11b4b",
					"request_id":"f77e8c6e-4852-47ec-b52a-786f1816a617"
				  }
				]
			  }
			`},
			expect: "{\"msg\":\"Data send error to firehose as well as Dead letter queue send error\",\"success\":false}",
			err:    main.ErrorDeadLetterQueue,
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.context, test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}
