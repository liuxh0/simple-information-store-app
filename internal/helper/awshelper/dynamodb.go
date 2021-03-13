package awshelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetDynamoDbClient returns a DynamoDB client.
func GetDynamoDbClient(endpoint string) *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithEndpoint(endpoint)
	dynamoDbClient := dynamodb.New(sess, config)
	return dynamoDbClient
}
