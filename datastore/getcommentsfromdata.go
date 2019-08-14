package datastore

import (
	"context"
	"hypeman/comment"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const comments = "Comments"

//Retrieves Username, Body, Date, Score from the database!
func (ds *DataStore) GetCommentsFromDB(videoname string) (*[]*comment.Comment, error) {

	threadComments := []*comment.Comment{}

	filter := bson.D{primitive.E{Key: "videoname", Value: videoname}}

	database := ds.Client.Database("Hypeman")
	cur, err := database.Collection("Comments").Find(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem comment.Comment
		err := cur.Decode(&elem)

		if err != nil {
			log.Fatal(err)
		}

		filter = bson.D{{"videoname", videoname}, {"commentid", elem.CommentID}}

		upvotes, _ := database.Collection("CommentLikes").CountDocuments(context.TODO(), filter)
		downvotes, _ := database.Collection("CommentDislikes").CountDocuments(context.TODO(), filter)

		elem.Score = int(upvotes*2 - downvotes)

		threadComments = append(threadComments, &elem)
	}

	return &threadComments, nil
}
