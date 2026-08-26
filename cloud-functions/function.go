package function

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/TheLuQ/eChart-backend/firestore"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
)

var sheetConnector *firestore.SheetConnector

func init() {
	dbConnector, initError := firestore.New("(default)")
	if initError != nil {
		println("Error initializing Firestore DB connector: " + initError.Error())
	}
	sheetConnector = &firestore.SheetConnector{Db: dbConnector, CollectionName: os.Getenv("DB")}
	functions.CloudEvent("AddEvent", AddEvent)
	functions.CloudEvent("RemoveEvent", RemoveEvent)
	functions.CloudEvent("MetadataUpdateEvent", ChangeMetadataEvent)
}

func ChangeMetadataEvent(ctx context.Context, e event.Event) error {
	sheet, err := parseEvent(e)
	if err != nil {
		return err
	}
	err = sheetConnector.UpsertSheet(sheet)
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
	println("Parsed sheet: " + sheet.Id + " with title: " + sheet.Title + " and name: " + sheet.Name)

	err = sheetConnector.AddSheet(sheet)
	if err != nil {
		println("Error saving sheet to database: " + err.Error())
		return err
	}
	println("File [" + sheet.FileName + "] added to database: " + sheet.Id)
	return nil
}

func RemoveEvent(ctx context.Context, e event.Event) error {
	sheet, err := parseEvent(e)
	if err != nil {
		return err
	}
	err = sheetConnector.RemoveSheet(sheet)
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
		return nil, err
	}
	configuredBucket := strings.TrimSpace(os.Getenv("SHEET_STORAGE_BUCKET"))
	if configuredBucket != "" && sth.Bucket != configuredBucket {
		return nil, fmt.Errorf("ignoring object from unexpected bucket: got %q expected %q", sth.Bucket, configuredBucket)
	}
	sheets, err := firestore.ToSheet(sth.Name, sth.Metadata)
	if err != nil {
		println("Error creating sheet from path: " + err.Error())
		return nil, err
	}
	return sheets, nil
}
