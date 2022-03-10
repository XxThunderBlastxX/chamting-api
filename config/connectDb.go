package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func ConnectDb() {
	//Get MongoUri from .env
	var uri = Env("MONGO_URI")

	//Passed the MongoUri
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Connecting to Mongo
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection established with MongoDB 🙌")
	}
	defer client.Disconnect(ctx)

	//Ping to the MongoServer
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
}
