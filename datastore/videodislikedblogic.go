package datastore

import (
	"context"
	"hypeman/metadata"
	"hypeman/opinions"

	"go.mongodb.org/mongo-driver/bson"
)

//if dislike == true, add if it doesn't exist. if false, remove if it does exist!
func (ds *DataStore) PushDislikeToDB(videoname string, username string, dislike bool) {
	coll := ds.Client.Database("Hypeman").Collection("VideoDislikes")

	item := opinions.VideoRating{Username: username, Videoname: videoname}

	if dislike { //if dislike is true, add it if it doesn't exist!
		filter := bson.D{{"videoname", videoname}, {"username", username}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count == 0 {
			coll.InsertOne(context.TODO(), item)

			//Update the cache
			md, err := ds.TimeCache.GetMetadata(videoname)

			if err == nil {
				ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "DISLIKES", Delta: 1}
			}
		}
	} else { //if dislike is false, delete it if it exists

		filter := bson.D{{"videoname", videoname}, {"username", username}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count > 0 {
			coll.DeleteOne(context.TODO(), filter)

			//Update the cache
			md, err := ds.TimeCache.GetMetadata(videoname)

			if err == nil {
				ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "DISLIKES", Delta: -1}
			}
		}
	}
}
