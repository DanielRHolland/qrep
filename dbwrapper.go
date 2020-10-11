package main

import (
	"context"
	"fmt"
	"log"
	"time"

	//   "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//   "go.mongodb.org/mongo-driver/mongo/readpref"
)

const mongodbURI = "mongodb://localhost:27017"
const dbname = "qrepdb"

func connectdb() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	item := trackedItem{"defaultName", nil}
	insertTrackedItem(ctx, client, item)

}

func insertTrackedItem(ctx context.Context, client *mongo.Client, item trackedItem) {
	collection := client.Database(dbname).Collection("tracked_items")
	insertResult, err := collection.InsertOne(ctx, item)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted trackedItem with ID:", insertResult.InsertedID)
}
