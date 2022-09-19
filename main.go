package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	var uri string

	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	//begin findOne
	coll := client.Database("users").Collection("subscriptions")

	doc := bson.D{{"chat_id", "5708402489"}, {"lat", 29.08}, {"lon", 48.05}}

	result, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.InsertedID, "id of inserted doc")

	select {}

}
