package database

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	mongo *mongo.Client
	redis *redis.Client
}

func NewDatabase(mongo *mongo.Client, redis *redis.Client) *Database {
	return &Database{
		mongo: mongo,
		redis: redis,
	}
}

//first check to ensure that user credentials are A-okay!
func (db *Database) ValidateUserIdAndPassword(userid int64, password string) error {
	collection := db.mongo.Database("Hypepeople").Collection("Users")
	filter := bson.M{"userid": userid, "password": password}

	count, err := collection.CountDocuments(context.TODO(), filter)

	if err != nil || count == 0 {
		return errors.New("Invalid userid or password")
	}

	return nil

}
