package db

import (
	"context"
	"log"
	"project-mgmt-go/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mongoInstance MongoInstance

func ConnectDB() {
	mongoURI := config.Env("MONGO_URI")
	dbName := config.Env("DB_NAME")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)

	mongoInstance = MongoInstance{
		Client: client,
		Db:     db,
	}
}

func Collection(collectionName string) *mongo.Collection {
	collection := mongoInstance.Db.Collection(collectionName)
	return collection
}
