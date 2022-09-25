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
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DbClient struct {
	Coll   *mongo.Collection
	Client *mongo.Client
}

type Subscription struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ChatId     int                `bson:"chat_id"`
	Lon        float64            `bson:"lon,required"`
	Lat        float64            `bson:"lat,required"`
	UpdateTime int                `bson:"update_time"`
}

func NewDbClient(cfg Config, log *zerolog.Logger) *DbClient {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DbTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DbHost))
	if err != nil {
		log.Fatal()
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal().Err(err).Send()
	}
	coll := client.Database(cfg.Db).Collection(cfg.Collection)

	return &DbClient{
		Client: client,
		Coll:   coll,
	}

}
func (db *DbClient) Disconnect(ctx context.Context) error {
	return db.Client.Disconnect(ctx)

}

func (db *DbClient) InsertOne(ctx context.Context, ins Subscription) (primitive.ObjectID, error) {
	result, err := db.Coll.InsertOne(ctx, ins)

	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("couldn't insert object: %w", err)
	}

	id, _ := result.InsertedID.(primitive.ObjectID)

	return id, nil

}

func (db *DbClient) UpsertOne(ctx context.Context, ins Subscription) (*mongo.UpdateResult, error) {

	filter := bson.D{{"chat_id", ins.ChatId}}
	update := bson.D{{"$set", bson.D{{"chat_id", ins.ChatId}, {"lat", ins.Lat}, {"lon", ins.Lon},
		{"update_time", ins.UpdateTime}}}}
	opts := options.Update().SetUpsert(true)
	result, err := db.Coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("couldn't update subescription: %w", err)
	}
	return result, nil

}

func (db *DbClient) FindOne(ctx context.Context, filter bson.D) (Subscription, error) {

	var res Subscription

	err := db.Coll.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return Subscription{}, fmt.Errorf("no subscription with this id: %w", err)

	} else if err != nil {
		return Subscription{}, err
	}

	return res, nil

}

func (db *DbClient) Find(ctx context.Context, filter bson.D) ([]Subscription, error) {

	var res []Subscription

	cursor, err := db.Coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding elements: %w", err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &res); err != nil {
		return nil, fmt.Errorf("error exctrating elements from cursor: %w", err)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error given when finding elements: %w", err)
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
	return nil
}
