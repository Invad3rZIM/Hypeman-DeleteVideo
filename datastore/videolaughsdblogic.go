package datastore

import (
	"context"
	"hypeman/metadata"
	"hypeman/opinions"

	"go.mongodb.org/mongo-driver/bson"
)

//if laugh == true, add the laugh if it doesn't exist. if false, remove if it does exist!
func (ds *DataStore) PushLaughToDB(videoname string, username string, laugh bool) {
	coll := ds.Client.Database("Hypeman").Collection("VideoLaughs")

	item := opinions.VideoRating{Username: username, Videoname: videoname}

	if laugh { //if laugh is true, add it if it doesn't exist!
		filter := bson.D{{"videoname", videoname}, {"username", username}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count == 0 {
			coll.InsertOne(context.TODO(), item)

			//Update the cache
			md, err := ds.TimeCache.GetMetadata(videoname)

			if err == nil {
				ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "LAUGHS", Delta: 1}
			}
		}
	} else { //if laugh is false, delete it if it exists

		filter := bson.D{{"videoname", videoname}, {"username", username}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count > 0 {
			coll.DeleteOne(context.TODO(), filter)
		}

		//Update the cache
		md, err := ds.TimeCache.GetMetadata(videoname)

		if err == nil {
			ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "LAUGHS", Delta: -1}
		}
	}
}
