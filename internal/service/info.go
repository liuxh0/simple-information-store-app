package service

import (
	"fmt"
	"simple-information-store-app/internal/env"
	"simple-information-store-app/internal/helper"
	"simple-information-store-app/internal/helper/awshelper"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Info presents an info item.
type Info struct {
	ID    string
	Value string
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../servicefakes . InfoCreator
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../servicefakes . InfoGetter
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../servicefakes . InfoUpdater

type InfoCreator interface {
	// CreateInfo creates an info in the database.
	// ValueTooLongError is returned if the value length exceeds the limit.
	CreateInfo(id, value string) (Info, error)
}

type InfoGetter interface {
	// GetInfo returns the info with the given id.
	// InfoNotFoundError is returned if the info does not exist.
	GetInfo(id string) (Info, error)
}

type InfoUpdater interface {
	// UpdateInfo updates an existing info.
	// ValueTooLongError is returned if the value length exceeds the limit.
	// InfoNotFoundError is returned if the info does not exist.
	UpdateInfo(id, newValue string) (Info, error)
}

type InfoService interface {
	InfoCreator
	InfoGetter
	InfoUpdater

	// DeleteInfo deletes an existing info.
	DeleteInfo(id string) error
}

type infoService struct{}

func NewInfoService() InfoService {
	return infoService{}
}

// ValueTooLongError indicates that the value length exceeds the limit.
type ValueTooLongError struct {
	AllowedLen int
	ActualLen  int
}

func (err ValueTooLongError) Error() string {
	return fmt.Sprintf("The length of the value is %d, however max. %d allowed.", err.ActualLen, err.AllowedLen)
}

// InfoNotFoundError indicates that the info with the given id does not exist.
type InfoNotFoundError struct {
	InfoID string
}

func (err InfoNotFoundError) Error() string {
	return fmt.Sprintf("Info with id %s does not exist.", err.InfoID)
}

func (_ infoService) CreateInfo(id, value string) (Info, error) {
	if err := checkValueLen(value); err != nil {
		return Info{}, *err
	}

	dynamoDbClient := awshelper.GetDynamoDbClient(env.GetDynamoDbEndpoint())
	valueTableName := env.GetValueTableName()
	_, err := dynamoDbClient.PutItem(&dynamodb.PutItemInput{
		TableName:           &valueTableName,
		ConditionExpression: helper.StringPtr("attribute_not_exists(Id)"),
		Item: map[string]*dynamodb.AttributeValue{
			"Id":    {S: &id},
			"Value": {S: &value},
		},
	})

	if err != nil {
		return Info{}, err
	}

	return Info{
		ID:    id,
		Value: value,
	}, nil
}

func (_ infoService) GetInfo(id string) (Info, error) {
	dynamoDbClient := awshelper.GetDynamoDbClient(env.GetDynamoDbEndpoint())
	valueTableName := env.GetValueTableName()
	result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &valueTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: &id},
		},
	})

	if err != nil {
		return Info{}, err
	}

	if result.Item == nil {
		return Info{}, InfoNotFoundError{
			InfoID: id,
		}
	}

	return Info{
		ID:    id,
		Value: *result.Item["Value"].S,
	}, nil
}

func (_ infoService) UpdateInfo(id, newValue string) (Info, error) {
	if err := checkValueLen(newValue); err != nil {
		return Info{}, *err
	}

	dynamoDbClient := awshelper.GetDynamoDbClient(env.GetDynamoDbEndpoint())
	valueTableName := env.GetValueTableName()

	// First check if the id exists
	result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &valueTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: &id},
		},
	})

	if err != nil {
		return Info{}, err
	}

	if result.Item == nil {
		return Info{}, InfoNotFoundError{
			InfoID: id,
		}
	}

	// If the id exists, update the value
	_, err = dynamoDbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: &valueTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: &id},
		},
		ConditionExpression: helper.StringPtr("attribute_exists(Id)"),
		UpdateExpression:    helper.StringPtr("set #Value = :value"),
		ExpressionAttributeNames: map[string]*string{
			"#Value": helper.StringPtr("Value"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":value": {S: &newValue},
		},
	})

	if err != nil {
		return Info{}, err
	}

	return Info{
		ID:    id,
		Value: newValue,
	}, nil
}

func (_ infoService) DeleteInfo(id string) error {
	dynamoDbClient := awshelper.GetDynamoDbClient(env.GetDynamoDbEndpoint())
	valueTableName := env.GetValueTableName()
	_, err := dynamoDbClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &valueTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: &id},
		},
	})

	return err
}

func checkValueLen(value string) *ValueTooLongError {
	if l := len(value); l > ValueMaxLen {
		return &ValueTooLongError{
			AllowedLen: ValueMaxLen,
			ActualLen:  l,
		}
	}

	return nil
}
