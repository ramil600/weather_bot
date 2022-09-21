package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbClient struct {
	Coll *mongo.Collection
}

func NewDbClient(cfg Config, log *zerolog.Logger) *DbClient {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DbTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DbHost))
	if err != nil {
		log.Fatal()
	}
	coll := client.Database(cfg.Db).Collection(cfg.Collection)

	return &DbClient{
		Coll: coll,
	}

}

func (db *DbClient) InsertOne(ctx context.Context, ins Subscription) (primitive.ObjectID, error) {
	result, err := db.Coll.InsertOne(ctx, ins)

	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("couldn't insert object: %w", err)
	}

	id, _ := result.InsertedID.(primitive.ObjectID)

	return id, nil

}

func (db *DbClient) FindOne(ctx context.Context, id primitive.ObjectID) (Subscription, error) {

	var res Subscription
	filter := bson.D{{"_id", id}}

	err := db.Coll.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return Subscription{}, fmt.Errorf("no subscription with this id: %w", err)

	} else if err != nil {
		return Subscription{}, err
	}

	return res, nil

}

func (db *DbClient) UpdateByID(ctx context.Context, id primitive.ObjectID, upd Subscription) (int, error) {

	updBson := bson.D{{"lat", upd.Lat}, {"lon", upd.Lon}}
	updres, err := db.Coll.UpdateByID(ctx, id,
		bson.D{{"$set", updBson}})
	if updres.MatchedCount == 0 {
		return 0, fmt.Errorf("could not match the subscription for update")
	}
	if err != nil {
		return 0, fmt.Errorf("UpdatById returned an error %w", err)
	}
	return int(updres.ModifiedCount), nil

}

func (db *DbClient) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}
	del, err := db.Coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete subscription: %w", err)
	}
	if del.DeletedCount == 0 {
		return fmt.Errorf("no record to delete")

	}
}
