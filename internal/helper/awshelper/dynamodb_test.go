package awshelper_test

import (
	"simple-information-store-app/internal/helper/awshelper"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDynamoDbClient()", func() {
	const testEndpoint = "http://test-endpoint.com/dynamodb"
	var dynamoDbClient *dynamodb.DynamoDB

	BeforeEach(func() {
		dynamoDbClient = awshelper.GetDynamoDbClient(testEndpoint)
	})

	It("should return a DynamoDB client with correct endpoint", func() {
		Expect(dynamoDbClient.Endpoint).To(Equal(testEndpoint))
	})
})
