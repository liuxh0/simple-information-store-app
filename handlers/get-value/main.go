package main

import (
	"fmt"

	"simple-information-store-app/internal/awshelper"
	"simple-information-store-app/internal/env"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	dynamoDbClient := awshelper.GetDynamoDbClient()
	tableName := env.GetValueTableName()
	result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: &id,
			},
		},
	})

	if err != nil {
		fmt.Printf("Error when retrieving item: %s\n", err.Error())
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	}

	value := result.Item["Value"].S
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       *value,
	}, nil
}

func main() {
	lambda.Start(handler)
}
