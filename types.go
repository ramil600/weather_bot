package main

type Subscription struct {
	ChatId string  `bson:"chat_id"`
	Lon    float64 `bson:"lat,required"`
	Lat    float64 `bson:"lon,required"`
}
