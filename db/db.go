package db

import (
	"context"
	"errors"
	"github.com/credo-science/credo-metadata-server/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var ErrNotFound = errors.New("not found")

var database *mongo.Database

func connect() {
	url, exists := os.LookupEnv("MONGO_URL")
	if !exists {
		url = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	database = client.Database("credo-metadata")
}

func init() {
	connect()
}

func getCollection(et event.Type) *mongo.Collection {
	return database.Collection(et.String())
}

func GetEventMetadata(et event.Type, id string) (event.Metadata, error) {
	col := getCollection(et)

	filter := bson.M{"id": id}
	projection := bson.M{"_id": 0, "id": 0}

	result := event.Metadata{}

	err := col.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return result, nil
}

func SetEventMetadata(et event.Type, id string, m event.Metadata) error {
	col := getCollection(et)

	filter := bson.M{"id": id}

	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": m}, options.Update().SetUpsert(true))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		} else {
			return err
		}
	}

	return nil
}
