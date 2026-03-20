package function

import (
	"context"
	"encoding/json"
	"path"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/TheLuQ/eChart-backend/firestore"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
)

var dbConnector *firestore.FireDb

func init() {
	var initError error
	dbConnector, initError = firestore.New("(default)", os.Getenv("DB"))
	if initError != nil {
		println("Error initializing Firestore DB connector: " + initError.Error())
	}
	functions.CloudEvent("AddEvent", AddEvent)
	functions.CloudEvent("RemoveEvent", RemoveEvent)
	functions.CloudEvent("MetadataUpdateEvent", ChangeMetadataEvent)
}

func ChangeMetadataEvent(ctx context.Context, e event.Event) error {
	sheet, err := parseEvent(e)
	if err != nil {
		return err
	}
	err = dbConnector.UpdateDocumentWithParentPath(path.Dir(sheet.Id), firestore.UpsertSheetFn(sheet))
	if err != nil {
		println("Error saving sheet to database: " + err.Error())
		return err
	}
	return nil
}

func AddEvent(ctx context.Context, e event.Event) error {
	sheet, err := parseEvent(e)
	if err != nil {
		return err
	}
	parentPath := path.Dir(sheet.Id)

	err = dbConnector.UpdateDocumentWithParentPath(parentPath, firestore.AddSheetFn(sheet, parentPath))
	if err != nil {
		println("Error saving sheet to database: " + err.Error())
		return err
	}
	return nil
}

func RemoveEvent(ctx context.Context, e event.Event) error {
	sheet, err := parseEvent(e)
	if err != nil {
		return err
	}
	err = dbConnector.UpdateDocumentWithParentPath(path.Dir(sheet.Id), firestore.RemoveSheetFn(sheet))
	if err != nil {
		println("Error saving sheet to database: " + err.Error())
		return err
	}
	return nil
}

func parseEvent(e event.Event) (*firestore.Sheet, error) {
	var sth storagedata.StorageObjectData
	if err := json.Unmarshal(e.Data(), &sth); err != nil {
		println("Error unmarshaling data: " + err.Error())
	}
	sheets, err := firestore.ToSheet(sth.Name, sth.Metadata)
	if err != nil {
		println("Error creating sheet from path: " + err.Error())
		return nil, err
	}
	return sheets, nil
}
