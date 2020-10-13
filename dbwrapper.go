package main

import (
	"context"
	"fmt"
	"log"
	"time"

	   "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//   "go.mongodb.org/mongo-driver/mongo/readpref"
)

const mongodbURI = "mongodb://localhost:27017"
const dbname = "qrepdb"
const itemsCollection = "tracked_items" 


func connectdb() (context.Context, *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
//	defer client.Disconnect(ctx)
        return ctx, client

}

func insertTrackedItem(ctx context.Context, client *mongo.Client, item trackedItem) string {
	collection := client.Database(dbname).Collection(itemsCollection)
	insertResult, err := collection.InsertOne(ctx, item)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted trackedItem with ID:", insertResult.InsertedID)
        return (insertResult).InsertedID.(primitive.ObjectID).Hex()
}

func insertItem(item trackedItem) string {
        ctx, client := connectdb()
        defer client.Disconnect(ctx)
	return insertTrackedItem(ctx, client, item)
}

func getItem(id string) (trackedItem, error) {
        objectId, _ := primitive.ObjectIDFromHex(id)
        ctx, client := connectdb()
        defer client.Disconnect(ctx)
        collection := client.Database(dbname).Collection(itemsCollection)
        filter:= bson.M{"_id": objectId}
        var item trackedItem
        err := collection.FindOne(ctx, filter).Decode(&item)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println("Found item with name ", item.Name)

        return item, err
}

func addIssueToItem(issue string, id string) {
        objectId, _ := primitive.ObjectIDFromHex(id)
        ctx, client := connectdb()
        defer client.Disconnect(ctx)
        collection := client.Database(dbname).Collection(itemsCollection)
        filter := bson.M{"_id": objectId}
        update := bson.M{"$push":bson.M{"issues": issue}}
        _, err := collection.UpdateOne(ctx,filter,update)
        if err !=nil{
            log.Fatal(err)
        }else{
        fmt.Println("success")
        }



}
