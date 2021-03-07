package main

import (
	"encoding/json"
	"fmt"

	"simple-information-store-app/internal/env"
	"simple-information-store-app/internal/helper"
	"simple-information-store-app/internal/helper/awshelper"
	"simple-information-store-app/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body

	if bodyLen := len(body); bodyLen > service.ValueMaxLen {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Actual length %d is greater than allowed length %d.", bodyLen, service.ValueMaxLen),
		}, nil
	}

	// Generate an Id
	id := uuid.New().String()

	dynamoDbClient := awshelper.GetDynamoDbClient()
	tableName := env.GetValueTableName()
	_, err := dynamoDbClient.PutItem(&dynamodb.PutItemInput{
		TableName:           &tableName,
		ConditionExpression: helper.PointerToString("attribute_not_exists(Id)"),
		Item: map[string]*dynamodb.AttributeValue{
			"Id":    {S: &id},
			"Value": {S: &body},
		},
	})

	if err != nil {
		fmt.Printf("Error when putting item: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	responseBody := map[string]string{
		"id": id,
	}
	responseBodyBytes, _ := json.Marshal(responseBody)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBodyBytes),
	}, nil
}

func main() {
	lambda.Start(handler)
}
