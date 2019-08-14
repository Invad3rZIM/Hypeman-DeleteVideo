package datastore

import (
	"context"
	"hypeman/metadata"
	"time"
)

//Uploads video metadata to the database and returns a pointer to said data, otherwise an error
func (ds *DataStore) UploadMetadataToDB(username string, songname string, tags []string) (*metadata.Metadata, error) {
	item := metadata.Metadata{
		Videoname: username + "*" + songname,
		Artist:    username,
		Songname:  songname,
		Tags:      tags,
		Score:     0,
		Likes:     0,
		Dislikes:  0,
		Laughs:    0,
		Comments:  0,
		Bucket:    "TODAY", //TODAY, THISWEEK, THISMONTH, THISYEAR, ALLTIME Are your arbitrary options...
		Views:     0,
		Hash:      0,
		Date:      time.Now().Unix(),
	}

	videos := ds.Client.Database("Hypeman").Collection("Metadata")
	_, err := videos.InsertOne(context.TODO(), item)

	if err != nil {
		return nil, err
	}

	ds.TimeCache.AddToCache(&item, "")

	return &item, nil
}
