package datastore

import (
	"context"
	"hypeman/metadata"
	"hypeman/opinions"

	"go.mongodb.org/mongo-driver/bson"
)

//if laugh == true, add the laugh if it doesn't exist. if false, remove if it does exist!
func (ds *DataStore) PushLikeToDB(videoname string, username string, like bool) {
	coll := ds.Client.Database("Hypeman").Collection("VideoLikes")

	item := opinions.VideoRating{Username: username, Videoname: videoname}

	if like { //if like is true, add it if it doesn't exist!
		filter := bson.D{{"videoname", videoname}, {"username", username}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count == 0 {
			coll.InsertOne(context.TODO(), item)

			//Update the cache
			md, err := ds.TimeCache.GetMetadata(videoname)

			if err == nil {
				ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "LIKES", Delta: 1}
			}
		}
	} else { //if like is false, delete it if it exists

		filter := bson.D{{"videoname", videoname}, {"username", username}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count > 0 {
			coll.DeleteOne(context.TODO(), filter)

			//Update the cache
			md, err := ds.TimeCache.GetMetadata(videoname)

			if err == nil {
				ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "LIKES", Delta: -1}
			}
		}
	}
}
