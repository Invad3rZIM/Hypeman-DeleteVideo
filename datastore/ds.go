package datastore

import (
	"hypeman/cache"

	"go.mongodb.org/mongo-driver/mongo"
)

type DataStore struct {
	Client    *mongo.Client
	TimeCache *cache.TimeCache
}

/*
const videos = "Videos"
const today = "Today"

//Retrieves Username, Body, Date, Score from the database!
func (ds *DataStore) AllVideosTodayDB() (*[]*Metadata, error) {
	filt := expression.Name("ThreadID").Equal(expression.Value(threadname))

	proj := expression.NamesList(expression.Name("Username"), expression.Name("Body"), expression.Name("Date"), expression.Name("Score"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// snippet-end:[dynamodb.go.scan_items.expr]

	// snippet-start:[dynamodb.go.scan_items.call]
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(comments),
	}

	// Make the DynamoDB Query API call
	result, err := ds.svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	threadComments := []*Comment{}

	for _, i := range result.Items {
		item := Comment{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		threadComments = append(threadComments, &item)
	}

	return &threadComments, nil
}
*/
