package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// type KinesisFirehoseEventRecordData struct {
// 	CustomerId string `json:"customerId"`
// }
// type JsonEvents []JsonEvent

type JsonEvent struct {
	AppID     string `json:"app_id"`
	DeviceID  string `json:"device_id"`
	RequestID string `json:"request_id"`
}

func handleRequest(evnt events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {

	var response events.KinesisFirehoseResponse

	for _, record := range evnt.Records {

		var transformedRecord events.KinesisFirehoseResponseRecord
		transformedRecord.RecordID = record.RecordID
		transformedRecord.Result = events.KinesisFirehoseTransformedStateOk

		var metaData events.KinesisFirehoseResponseRecordMetadata

		var eventRecordData JsonEvent

		partitionKeys := make(map[string]string)

		json.Unmarshal(record.Data, &eventRecordData)

		partitionKeys["app_id"] = eventRecordData.AppID

		metaData.PartitionKeys = partitionKeys

		transformedRecord.Metadata = metaData

		csv_data := fmt.Sprintf("%s,%s,%s\n", eventRecordData.AppID, eventRecordData.DeviceID, eventRecordData.RequestID)

		str := base64.StdEncoding.EncodeToString([]byte(csv_data))

		data, err := base64.StdEncoding.DecodeString(str)

		if err != nil {
			log.Fatal("error:", err)
		}

		record.Data = data

		transformedRecord.Data = record.Data

		response.Records = append(response.Records, transformedRecord)

	}
	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
