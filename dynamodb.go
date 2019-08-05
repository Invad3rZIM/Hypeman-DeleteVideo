// snippet-comment:[These are tags for the AWS doc team's sample catalog. Do not remove.]
// snippet-sourceauthor:[Doug-AWS]
// snippet-sourcedescription:[DynamoDBReadItem.go gets an item from an Amazon DynamoDB table.]
// snippet-keyword:[Amazon DynamoDB]
// snippet-keyword:[GetItem function]
// snippet-keyword:[Go]
// snippet-service:[dynamodb]
// snippet-keyword:[Code Sample]
// snippet-sourcetype:[full-example]
// snippet-sourcedate:[2019-03-19]
/*
   Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at
    http://aws.amazon.com/apache2.0/
   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/
// snippet-start:[dynamodb.go.read_item]
package main

// snippet-start:[dynamodb.go.read_item.imports]

// snippet-end:[dynamodb.go.read_item.imports]

// snippet-start:[dynamodb.go.read_item.struct]
// Create struct to hold info about new item
type Item struct {
	Username, First, Last string
	Followers, Following  int
}

// snippet-end:[dynamodb.go.read_item.struct]
/*
func main2(svc *dynamodb.DynamoDB) string {
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
}
*/
