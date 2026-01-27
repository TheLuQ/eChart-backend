package main

import (
	"context"

	function "github.com/TheLuQ/eChart-backend"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		panic("Failed to create client, " + err.Error())
	}

	println("Starting serving cloud event...")
	if err = c.StartReceiver(context.Background(), function.AddEvent); err != nil {
		panic("Failed to start receiver " + err.Error())
	}

}
