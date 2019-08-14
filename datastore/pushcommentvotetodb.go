package datastore

import (
	"context"
	"hypeman/opinions"

	"go.mongodb.org/mongo-driver/bson"
)

//if laugh == true, add the laugh if it doesn't exist. if false, remove if it does exist!
func (ds *DataStore) PushCommentVoteToDB(col string, videoname string, username string, commentid string, vote int) {
	coll := ds.Client.Database("Hypeman").Collection(col)

	item := opinions.CommentRating{Videoname: videoname, CommentID: commentid, Username: username}

	if vote == 1 { //if like is true, add it if it doesn't exist!
		filter := bson.D{{"videoname", videoname}, {"username", username}, {"commentid", commentid}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count == 0 {
			coll.InsertOne(context.TODO(), item)

			//Update the cache
			//	md, err := ds.TimeCache.GetMetadata(videoname)

			//	if err == nil {
			//	ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "LIKES", Delta: 1}
			//	}
		}
	} else if vote == 0 { //if like is false, delete it if it exists

		filter := bson.D{{"videoname", videoname}, {"username", username}, {"commentid", commentid}}

		count, _ := coll.CountDocuments(context.TODO(), filter)

		if count > 0 {
			coll.DeleteOne(context.TODO(), filter)

			//Update the cache
			//		md, err := ds.TimeCache.GetMetadata(videoname)

			//		if err == nil {
			//		ds.TimeCache.Changes <- &metadata.DataChange{Videoname: md.Videoname, Category: "LIKES", Delta: -1}
			//		}
		}
	}
}
