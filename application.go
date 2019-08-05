package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/labstack/echo"
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

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, os.Getenv("user"))
	})

	e.GET("/database", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%+v", main2()))
	})

	e.Logger.Fatal(e.Start(":5000"))
}

//ec2-3-82-204-144.compute-1.amazonaws.com

func main2() string {
	return "XXX"
	//tableName := "Users"
	/*
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
				return ""
			}

		item := Item{}
		/*
			err = dynamodbattribute.UnmarshalMap(result.Item, &item)
			if err != nil {
				return fmt.Sprintf("Failed to unmarshal Record, %v", err)
			}

		return fmt.Sprintf("%+v", item)

		return fmt.Sprintf("%+v", svc.ClientInfo)*/
}
