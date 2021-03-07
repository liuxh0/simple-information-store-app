package awshelper

import (
	"simple-information-store-app/internal/env"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sharedDynamoDbClient *dynamodb.DynamoDB

func GetDynamoDbClient() *dynamodb.DynamoDB {
	if sharedDynamoDbClient == nil {
		sess := session.Must(session.NewSession())
		config := aws.NewConfig().WithEndpoint(env.GetDynamoDbEndpoint())
		sharedDynamoDbClient = dynamodb.New(sess, config)
	}
	return sharedDynamoDbClient
}
