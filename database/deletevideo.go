package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var videoProjection = bson.D{
	{"videoid", 1},
	{"previewid", 1},
	{"userid", 1},
	{"timebucket", 1},
	{"tags", 1},
}

//Metadata contains all the data that describes user uploaded videos
type VideoMetadata struct {
	VideoId, PreviewId string
	UserId             int64
	TimeBucket         string
	Tags               []string
}

//Retrieves video by ID from the DB; on failure return (nil, err) instead
func (db *Database) RetrieveVideo(filter *bson.M) (*VideoMetadata, error) {
	var result VideoMetadata

	cur, err := db.mongo.Database("Hypeman").Collection("Metadata").Find(
		context.Background(),
		*filter,
		options.Find().SetProjection(videoProjection))

	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		err := cur.Decode(&result)

		if err != nil {
			return nil, err
		}
	}

	if result.VideoId == "" {
		return nil, errors.New("Video not found in collection")
	}

	return &result, nil
}

//Drop all video ratings in every possible bucket
func (db *Database) purgeVideoRatings(filter *bson.M, ch chan bool) {

	colls := []string{"VideoLikes", "VideoDislikes", "VideoLaughs", "VideoViews"}

	for _, c := range colls {
		db.mongo.Database("Hypeman").Collection(c).DeleteMany(context.TODO(), *filter)
	}

	ch <- true
}

//Drop all video ratings in every possible bucket
func (db *Database) purgeVideoComments(filter *bson.M, ch chan bool) {
	colls := []string{"Comments", "CommentLikes"}

	for _, c := range colls {
		db.mongo.Database("Hypeman").Collection(c).DeleteMany(context.TODO(), *filter)
	}

	ch <- true

}

//Drop all video ratings in every possible bucket
func (db *Database) purgeRedis(record *VideoMetadata, ch chan bool) {
	cats := []string{"likes", "dislikes", "laughs", "views", "new"}

	//delete all instances in general cache
	for _, c := range cats {
		db.redis.ZRem(fmt.Sprintf("%s:%s", record.TimeBucket, c), record.VideoId)
	}

	//delete all instances in tags
	for _, t := range record.Tags {
		db.redis.ZRem(fmt.Sprintf("vidsbytag:%s", t))
	}

	db.redis.Del(fmt.Sprintf("timebucket:%s", record.VideoId))

	ch <- true
}

func (db *Database) DeleteVideo(userid int64, videoid string) (string, error) {
	filter := bson.M{"videoid": videoid}

	//grab video
	result, err := db.RetrieveVideo(&filter)
	if err != nil {
		return "", err
	}

	//OWNERSHIP CHECK - no ownership => admin check
	if result.UserId != userid {
		//ADMIN CHECK
		filter := bson.M{"adminid": userid}
		adminCollection := db.mongo.Database("Roles").Collection("Admins")
		count, _ := adminCollection.CountDocuments(context.TODO(), filter)

		if count == 0 {
			return "", errors.New("Insufficient Access Privileges")
		}
	}

	ch := make(chan bool, 3)

	go db.purgeRedis(result, ch)
	go db.purgeVideoComments(&filter, ch)
	go db.purgeVideoRatings(&filter, ch)

	<-ch
	<-ch
	<-ch

	result.Tags = nil
	result.TimeBucket = ""

	//insert to the droplist to be removed from S3 during CRON OPS (handled in different microservice)
	db.mongo.Database("Hypeman").Collection("Droplist").InsertOne(context.TODO(), &result)
	db.mongo.Database("Hypeman").Collection("Metadata").DeleteOne(context.TODO(), filter)

	return "deleted", nil
}
