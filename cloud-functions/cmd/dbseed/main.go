package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/TheLuQ/eChart-backend/firestore"
)

type MockData struct {
	Sheets []firestore.Sheet `json:"sheets"`
	Events []firestore.Event `json:"events"`
}

func main() {
	mockFilePath := os.Getenv("MOCK_DATA_FILE")
	dbConnector, initError := firestore.New("(default)")
	if initError != nil {
		log.Fatal("Failed to initialize Firestore DB connector: " + initError.Error())
	}
	eventConnector := &firestore.EventConnector{Db: dbConnector, CollectionName: os.Getenv("EVENTS_COLLECTION")}
	sheetConnector := &firestore.SheetConnector{Db: dbConnector, CollectionName: os.Getenv("DB")}

	mockData := getMockData(mockFilePath)
	for _, sheet := range mockData.Sheets {
		err := sheetConnector.AddSheet(&sheet)
		if err != nil {
			log.Fatal("Failed to add sheet: " + err.Error())
		} else {
			log.Println("Successfully added sheet: " + sheet.Title)
		}
	}
	for _, event := range mockData.Events {
		err := eventConnector.AddEvent(&event)
		if err != nil {
			log.Fatal("Failed to add event: " + err.Error())
		} else {
			log.Println("Successfully added event: " + event.Title)
		}
	}
}

func getMockData(filePath string) MockData {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Failed to open mock data file: " + err.Error())
	}
	defer file.Close()

	var mockData MockData
	if err := json.NewDecoder(file).Decode(&mockData); err != nil {
		log.Fatal("Failed to decode mock data: " + err.Error())
	}

	return mockData
}
