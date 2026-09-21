package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TheLuQ/eChart-backend/api/router"
	"github.com/TheLuQ/eChart-backend/firestore"
)

func main() {
	if emulatorHost := os.Getenv("DATASTORE_EMULATOR_HOST"); emulatorHost != "" {
		log.Println("Using Datastore emulator at " + emulatorHost)
	}
	dbConnector, initError := firestore.New("(default)")
	if initError != nil {
		log.Fatal("Failed to initialize Firestore DB connector: " + initError.Error())
	}
	eventConnector := &firestore.EventConnector{Db: dbConnector, CollectionName: os.Getenv("EVENTS_COLLECTION")}
	sheetConnector := &firestore.SheetConnector{Db: dbConnector, CollectionName: os.Getenv("DB")}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, router.NewRouter(sheetConnector, eventConnector)); err != nil {
		log.Fatal(err)
	}
}
