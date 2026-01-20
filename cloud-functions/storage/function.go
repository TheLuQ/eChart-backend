package function

import (
	"context"
	"encoding/json"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/TheLuQ/eChart-backend/firestore"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
)

var dbConnector firestore.InterfaceFirestoreDB

func init() {
	var initError error
	dbConnector, initError = firestore.New("(default)", "pupu")
	if initError != nil {
		println("Error initializing Firestore DB connector: " + initError.Error())
	}
	functions.CloudEvent("CloudEventFunc", CloudEventFunc)
}

func CloudEventFunc(ctx context.Context, e event.Event) error {
	var sth storagedata.StorageObjectData
	if err := json.Unmarshal(e.Data(), &sth); err != nil {
		println("Error unmarshaling data: " + err.Error())
	}
	sheet, err := firestore.ToSheet(sth.Name)
	if err != nil {
		println("Error creating sheet from path: " + err.Error())
		return err
	}

	err = dbConnector.UpdateSheetGroup(sheet.Id, firestore.Sheet{Instrument: sheet.Instrument, Id: sth.Name}, sheet.Title)
	if err != nil {
		println("Error saving sheet to database: " + err.Error())
		return err
	}
	return nil
}
