package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	//   "go.mongodb.org/mongo-driver/mongo/readpref"
         "encoding/base64"
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
        var idbytes [12]byte =  primitive.NewObjectID()
        var idbyteslice []byte = idbytes[:]
        item.Id = base64.StdEncoding.EncodeToString(idbyteslice)
	collection := client.Database(dbname).Collection(itemsCollection)
	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		log.Fatal(err)
	}
	return item.Id
}

func insertItem(item trackedItem) string {
	ctx, client := connectdb()
	defer client.Disconnect(ctx)
	return insertTrackedItem(ctx, client, item)
}

func getItem(id string) (trackedItem, error) {
	ctx, client := connectdb()
	defer client.Disconnect(ctx)
	collection := client.Database(dbname).Collection(itemsCollection)
	filter := bson.M{"_id": id}
	var item trackedItem
	err := collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found item with name ", item.Name)

	return item, err
}

func updateDbIssue(issue issueType) error {
	ctx, client := connectdb()
	defer client.Disconnect(ctx)
	collection := client.Database(dbname).Collection(itemsCollection)
	filter := bson.M{"issues._id": issue.Id}
	update := bson.M{"$set": bson.M{"issues.$": issue}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("success")
	}
	return err
}

func addIssueToItem(issue issueType, id string) error {
	ctx, client := connectdb()
	defer client.Disconnect(ctx)
	issue.Id = primitive.NewObjectID().Hex()
	collection := client.Database(dbname).Collection(itemsCollection)
	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"issues": issue}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("success")
	}
	return err
}

func getTrackedItems(maxcount int) ([]trackedItem, error) {
	ctx, client := connectdb()
	defer client.Disconnect(ctx)
	collection := client.Database(dbname).Collection(itemsCollection)
	var items []trackedItem
	cursor, err := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var item trackedItem
		if err = cursor.Decode(&item); err != nil {
			log.Fatal(err)
		} else {
			items = append(items, item)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return items, err
}

func searchTrackedItems(maxcount int, name string) ([]trackedItem, error) {
	ctx, client := connectdb()
	defer client.Disconnect(ctx)
	collection := client.Database(dbname).Collection(itemsCollection)
	filter := bson.M{"name": name}
	var items []trackedItem
	cursor, err := collection.Find(ctx, filter)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var item trackedItem
		if err = cursor.Decode(&item); err != nil {
			log.Fatal(err)
		} else {
			items = append(items, item)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return items, err
}

func getItemsFromIds(itemids []string) []trackedItem {
	var items []trackedItem
	for _, id := range itemids {
		item, _ := getItem(id) //TODO Catch errors
		items = append(items, item)
	}
	return items
}

func removeItemsFromDb(itemids []string) {
        log.Println("Place holder remove items by id here. Ids:", itemids)
}
