package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	svc = dbConnect()

}

func dbConnect() *dynamodb.DynamoDB {
	//dbconnect; env variables located in Elastic Beanstalk
	//CONFIGURATION -> SOFTWARE -> ENV PROPERTIES
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Println(err.Error())
		log.Fatal(err.Error())
	}

	return dynamodb.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))
}

func main() {
	h := Handler{}

	http.HandleFunc("/video/upload", h.s3UploadHandler)
	http.HandleFunc("/database/test", xx)

	http.ListenAndServe(":5000", nil)

}

//ec2-3-82-204-144.compute-1.amazonaws.com

func xx(w http.ResponseWriter, r *http.Request) {
	tableName := "Users"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String("kzimmer"),
			},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	json.NewEncoder(w).Encode(fmt.Sprintf("%+v", item))
	return

}

type Item struct {
	Username, First, Last string
	Followers, Following  int
}
