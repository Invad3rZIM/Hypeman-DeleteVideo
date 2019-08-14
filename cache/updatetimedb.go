package cache

import (
	"context"
	"hypeman/metadata"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Updates the time in the mongodb if it's an inaccurate bucket
func (tc *TimeCache) UpdateDBTime(md *metadata.Metadata) {
	filter := bson.D{primitive.E{Key: "videoname", Value: md.Videoname}}

	update := bson.M{"$set": bson.M{"videoname": md.Videoname, "bucket": md.Bucket}}

	tc.Client.Database("Hypeman").Collection("Metadata").UpdateOne(context.TODO(), filter, update)
}
