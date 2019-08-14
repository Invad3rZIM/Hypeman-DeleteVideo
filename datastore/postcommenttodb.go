package datastore

import (
	"context"
	"hypeman/comment"
)

//Posts comment to the MONGODB
func (ds *DataStore) PostCommentToDB(comment *comment.Comment) error {
	comments := ds.Client.Database("Hypeman").Collection("Comments")
	_, err := comments.InsertOne(context.TODO(), comment)

	return err
}
