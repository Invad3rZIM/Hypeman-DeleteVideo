package main

import (
	"context"
	"fmt"
	"hypeman/handler"
	"log"
	"net/http"

	"github.com/Sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

}

func main() {

	clientOptions := options.Client().ApplyURI("mongodb+srv://kzimmer:Testing123@cluster0-p7s3g.mongodb.net/test?retryWrites=true&w=majority")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	err = client.Ping(context.Background(), nil)

	h := handler.NewHandler(client)

	go h.DataStore.TimeCache.PerpetualSortRoutine(3)
	go h.DataStore.TimeCache.UpdateCacheRoutine()

	http.HandleFunc("/video/upload", h.S3UploadHandler)
	http.HandleFunc("/video/metadata", h.GetVideoMetadataHandler)
	http.HandleFunc("/video/comments/post", h.PostCommentHandler)

	http.HandleFunc("/video/comments/vote", h.CommentVoteHandler)
	http.HandleFunc("/video/like", h.VideoLikeHandler)
	http.HandleFunc("/video/dislike", h.VideoDislikeHandler)
	http.HandleFunc("/video/laugh", h.VideoLaughHandler)
	http.HandleFunc("/video/view", h.VideoViewHandler)
	http.HandleFunc("/video/most", h.GlobalMostXHandler)

	fmt.Println("Serving on localhost :5000")
	http.ListenAndServe(":5000", nil)
}
