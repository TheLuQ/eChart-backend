package function

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoHandler(w http.ResponseWriter, r *http.Request) {
	connString := os.Getenv("MONGO_CONNECTION_STRING")
	if connString == "" {
		http.Error(w, "MONGO_CONNECTION_STRING environment variable is not set", http.StatusInternalServerError)
		return
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connString))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	fmt.Fprintln(w, "Connected to MongoDB!")

	result, err := client.Database("test-collection")
	.Collection("pipi")
	.InsertOne(context.Background(), bson.M{"hello": "world", "number": 42})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting document: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Inserted document with ID: %v", result.InsertedID)
}
