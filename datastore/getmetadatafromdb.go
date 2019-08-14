package datastore

import (
	"context"
	"hypeman/metadata"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Retrieves all metadata sans individual comments from the Metadata table
func (ds *DataStore) GetVideoFromDB(videoname string) (*metadata.Metadata, error) {
	filter := bson.D{primitive.E{Key: "videoname", Value: videoname}}

	var result metadata.Metadata
	err := ds.Client.Database("Hypeman").Collection("Metadata").FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil, err
	}

	database := ds.Client.Database("Hypeman")

	filter = bson.D{{"videoname", videoname}}

	laughCount, _ := database.Collection("VideoLaughs").CountDocuments(context.TODO(), filter)
	viewCount, _ := database.Collection("VideoViews").CountDocuments(context.TODO(), filter)
	likeCount, _ := database.Collection("VideoLikes").CountDocuments(context.TODO(), filter)
	dislikeCount, _ := database.Collection("VideoDislikes").CountDocuments(context.TODO(), filter)

	result.Laughs = int(laughCount)
	result.Views = int(viewCount)
	result.Likes = int(likeCount)
	result.Dislikes = int(dislikeCount)

	hash := result.Laughs + result.Views + result.Likes + result.Dislikes

	if hash != result.Hash {
		update := bson.M{"$set": bson.M{"laughs": result.Laughs, "views": result.Views, "dislikes": result.Dislikes, "likes": result.Likes, "hash": hash}}

		//update values and hash in database if the hash is distinct!
		database.Collection("Metadata").UpdateOne(context.TODO(), filter, update)
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
