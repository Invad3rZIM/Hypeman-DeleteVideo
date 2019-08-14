package cache

import (
	"context"
	"hypeman/metadata"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Retrieves Metadata from the time buckets (TODAY, THISWEEK, THISMONTH, THISYEAR), as well as ranking data
func (tc *TimeCache) GetBucketVideosFromDB(bucket string) (*[]*metadata.Metadata, error) {
	projection := bson.D{
		{"videoname", 1},
		{"bucket", 1},
		{"date", 1},
	}

	coll := tc.Client.Database("Hypeman").Collection("Metadata")

	cur, err := coll.Find(
		context.Background(),
		bson.D{
			{"bucket", bucket},
		},
		options.Find().SetProjection(projection),
	)

	if err != nil {
		return nil, err
	}

	mds := []*metadata.Metadata{}

	//retrieve every item and push it to the array
	for cur.Next(context.TODO()) {
		var elem metadata.Metadata

		err := cur.Decode(&elem)

		if err != nil {
			log.Fatal(err)
		}

		elemP, err := tc.GetVideoFromDB(elem.Videoname)

		if err == nil {
			mds = append(mds, elemP)
		}
	}
	return &mds, nil
}
