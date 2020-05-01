package handler

import (
	"context"
	"errors"
	"fmt"
	"hypeman-deletevideo/constants"
	"hypeman-deletevideo/database"
	"os"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	Mongo    *mongo.Client
	Database *database.Database
}

//Attempts to init the handler with a redis and mongo configured from ENV variables
//if it doesn't work, try it again twice with a (2,5) second wait, and if that still fails,
//terminate process with error
func InitializeHandler() (*Handler, error) {
	var err error
	failureCount := 0

	var redis *redis.Client

	for redis == nil && failureCount < 3 {
		redis, err = connectRedis()

		if err != nil {
			if failureCount == 3 {
				return nil, errors.New("Could not connect to redis, please try again later")
			}
			failureCount += 1

			time.Sleep(2 * time.Second)
		}
	}

	failureCount = 0

	var mongo *mongo.Client

	for mongo == nil && failureCount < 3 {
		mongo, err = connectMongo()

		if err != nil {
			if failureCount == 3 {
				return nil, errors.New("Could not connect to mongo, please try again later")
			}
			failureCount += 1
			time.Sleep(2 * time.Second)
		}
	}

	return &Handler{
		Mongo:    mongo,
		Database: database.NewDatabase(mongo, redis),
	}, nil
}

func connectRedis() (*redis.Client, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := redis.Ping().Result()

	if err != nil {
		return nil, err
	}

	return redis, nil

}

func connectMongo() (*mongo.Client, error) {
	connectionString := fmt.Sprintf(constants.MONGO_TEMPLATE, os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASS"))

	clientOptions := options.Client().ApplyURI(connectionString)
	mongo, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	//test connection with a ping!
	err = mongo.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return mongo, nil
}
