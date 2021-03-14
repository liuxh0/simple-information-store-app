package integration_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"simple-information-store-app/internal/env"
	"simple-information-store-app/internal/helper/awshelper"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationSuite")
}

const (
	samHost          = "http://localhost:3000"
	dynamoDbEndpoint = "http://localhost:8000"
)

var (
	dynamoDbClient *dynamodb.DynamoDB
	valueTableName string
)

var _ = BeforeSuite(func() {
	var err error

	By("checking local server is running")
	_, err = http.Get(samHost)
	if err != nil {
		Fail("SAM local is not running.")
	}

	By("checking local DynamoDB is running")
	_, err = http.Get(dynamoDbEndpoint)
	if err != nil {
		Fail("Local DynamoDB is not running.")
	}

	By("setting environment variable AWS_SAM_LOCAL")
	os.Setenv("AWS_SAM_LOCAL", "true")
	os.Setenv("AWS_REGION", "eu-central-1")

	By("initializing variables")
	dynamoDbClient = awshelper.GetDynamoDbClient(dynamoDbEndpoint)
	valueTableName = env.GetValueTableName()
})

func readReadCloserOrDie(rc io.ReadCloser) string {
	bytes, err := ioutil.ReadAll(rc)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func convertJsonStringToMap(jsonString string) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func getStringFromJsonString(jsonString string, field string) (value string, ok bool) {
	jsonMap, err := convertJsonStringToMap(jsonString)
	if err != nil {
		return "", false
	}

	rawValue, ok := jsonMap[field]
	if !ok {
		return "", false
	}

	value, ok = rawValue.(string)
	return
}

func deleteItem(id string) error {
	_, err := dynamoDbClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &valueTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: &id},
		},
	})
	return err
}
