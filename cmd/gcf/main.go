package main

import (
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/TheLuQ/eChart-backend"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	println("Starting cloud function on port: " + port)
	if err := funcframework.Start(port); err != nil {
		panic("Failed to start cloud function: " + err.Error())
	}
}
