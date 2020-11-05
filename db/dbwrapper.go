package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	//   "go.mongodb.org/mongo-driver/mongo/readpref"
         "encoding/base64"
         qrep "github.com/DanielRHolland/qrep/models"
)

type MongoDbAccessor struct {
     client *mongo.Client
     ctx context.Context
     trackeditems *mongo.Collection
}


func NewMongoDbConnection()  *MongoDbAccessor {
        const mongodbURI = "mongodb://127.0.0.1:27017"
        const dbname = "qrepdb"
        const itemsCollection = "tracked_items"
	ctx  := context.Background()
	client, err := mongo.Connect(ctx,options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
	}
        trackeditems := client.Database(dbname).Collection(itemsCollection)
        return &MongoDbAccessor{client, ctx, trackeditems}
}

func (m *MongoDbAccessor) Disconnect() {m.client.Disconnect(m.ctx)}


func (m *MongoDbAccessor) InsertItem(item qrep.TrackedItemType) (string, error) {
        var idbytes [12]byte =  primitive.NewObjectID()
        var idbyteslice []byte = idbytes[:]
        item.Id = base64.StdEncoding.EncodeToString(idbyteslice)
        _, err := m.trackeditems.InsertOne(m.ctx, item)
        if err != nil {
                log.Fatal(err)
        }
	return item.Id, nil
}

func (m *MongoDbAccessor) GetItem(id string) (qrep.TrackedItemType, error) {
	filter := bson.M{"_id": id}
	var item qrep.TrackedItemType
        err := m.trackeditems.FindOne(m.ctx, filter).Decode(&item)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("Found item with name ", item.Name)
	return item, err
}

func (m *MongoDbAccessor) UpdateDbIssue(issue qrep.IssueType) error {
	filter := bson.M{"issues._id": issue.Id}
	update := bson.M{"$set": bson.M{"issues.$": issue}}
        var err error
        _, err = m.trackeditems.UpdateOne(m.ctx, filter, update)
        if err != nil {
                log.Fatal(err)
        } else {
                log.Println("success")
        }
	return err
}

func (m *MongoDbAccessor) AddIssueToItem(issue qrep.IssueType, id string) error {
	issue.Id = primitive.NewObjectID().Hex()
	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"issues": issue}}
        var err error
        _, err = m.trackeditems.UpdateOne(m.ctx, filter, update)
        if err != nil {
                log.Fatal(err)
        } else {
                fmt.Println("success")
        }
	return err
}

func (m *MongoDbAccessor) GetTrackedItems(maxcount int) ([]qrep.TrackedItemType, error) {
	var items []qrep.TrackedItemType
        var err error

        cursor, err := m.trackeditems.Find(m.ctx, bson.M{})
        defer cursor.Close(m.ctx)
        for cursor.Next(m.ctx) {
                var item qrep.TrackedItemType
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

func (m *MongoDbAccessor) SearchTrackedItems(maxcount int, name string) ([]qrep.TrackedItemType, error) {
	filter := bson.M{"name": name}
	var items []qrep.TrackedItemType
        var err error
        cursor, err := m.trackeditems.Find(m.ctx, filter)
        defer cursor.Close(m.ctx)
        for cursor.Next(m.ctx) {
                var item qrep.TrackedItemType
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

func (m *MongoDbAccessor) GetItemsFromIds(itemids []string) ([]qrep.TrackedItemType, error) {
	var items []qrep.TrackedItemType
	for _, id := range itemids {
                log.Println("Getting:",id)
		item, _ := m.GetItem(id) //TODO Catch errors
		items = append(items, item)
	}
	return items, nil
}

func (m *MongoDbAccessor) RemoveItemsFromDb(itemids []string) error{
        for _, id := range itemids {
            filter := bson.M{"_id": id}
            _, err := m.trackeditems.DeleteOne(m.ctx, filter)
            if err != nil {
                log.Fatal(err)
            }
        }
        return nil
}
