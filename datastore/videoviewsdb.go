package datastore

import (
	"context"
	"hypeman/metadata"
	"hypeman/opinions"

	"go.mongodb.org/mongo-driver/bson"
)

func (ds *DataStore) PushViewToDB(videoname string, username string) {
	coll := ds.Client.Database("Hypeman").Collection("VideoViews")

	item := opinions.VideoRating{Username: username, Videoname: videoname}

	filter := bson.D{{"videoname", videoname}, {"username", username}}

	count, _ := coll.CountDocuments(context.TODO(), filter)

	if count == 0 {
		coll.InsertOne(context.TODO(), item)

		//Update the cache
		md, err := ds.TimeCache.GetMetadata(videoname)

		if err == nil {
			ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "VIEWS", Delta: 1}
		}
	}
}
