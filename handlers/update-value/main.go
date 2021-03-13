package main

import (
	"fmt"

	"simple-information-store-app/internal/env"
	"simple-information-store-app/internal/helper"
	"simple-information-store-app/internal/helper/awshelper"
	"simple-information-store-app/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	value := request.Body

	if valueLen := len(value); valueLen > service.ValueMaxLen {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Actual length %d is greater than allowed length %d.", valueLen, service.ValueMaxLen),
		}, nil
	}

	dynamoDbClient := awshelper.GetDynamoDbClient(env.GetDynamoDbEndpoint())
	tableName := env.GetValueTableName()
	_, err := dynamoDbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: &id},
		},
		UpdateExpression: helper.StringPtr("set #Value = :value"),
		ExpressionAttributeNames: map[string]*string{
			"#Value": helper.StringPtr("Value"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":value": {S: &value},
		},
	})

	if err != nil {
		fmt.Printf("Error when updating item: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
