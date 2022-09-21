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

	//begin InsertOne
	coll := client.Database("users").Collection("subscriptions")

	subscription := Subscription{
		ChatId: "5708402489",
		Lat:    29.08,
		Lon:    48.05,
	}

	result, err := coll.InsertOne(context.TODO(), subscription)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.InsertedID, "id of inserted doc")

	//begin FindOne
	var res Subscription
	filter := bson.D{{"chat_id", "5708402489"}}
	err = coll.FindOne(context.TODO(), filter).Decode(&res)
	fmt.Println("Found one item", res)

	select {}

}
