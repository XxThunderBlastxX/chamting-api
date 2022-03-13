package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectDb() *mongo.Client {
	//Get MongoUri from .env
	//var uri = Env("MONGO_URI")
	uri := "mongodb://localhost:27017/chamting"

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

	return client
}
