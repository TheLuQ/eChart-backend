package function

import (
	"context"

	"encoding/json"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
)

func init() {
	functions.CloudEvent("CloudEventFunc", CloudEventFunc)
}

func CloudEventFunc(ctx context.Context, e event.Event) error {
	var sth storagedata.StorageObjectData
	if err := json.Unmarshal(e.Data(), &sth); err != nil {
		println("Error unmarshaling data: " + err.Error())
	}

	println("File in bucket:" + sth.Bucket + " of name: " + sth.Name + "was created")
	return nil
}
