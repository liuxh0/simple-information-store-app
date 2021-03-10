package awshelper_test

import (
	"simple-information-store-app/internal/env"
	"simple-information-store-app/internal/helper/awshelper"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDynamoDbClient()", func() {
	var dynamoDbClient *dynamodb.DynamoDB

	JustBeforeEach(func() {
		dynamoDbClient = awshelper.GetDynamoDbClient()
	})

	It("should return a DynamoDB client with correct endpoint", func() {
		endpoint := env.GetDynamoDbEndpoint()
		Expect(dynamoDbClient.Endpoint).To(Equal(endpoint))
	})

	When("calling it again", func() {
		var secondReturn *dynamodb.DynamoDB

		JustBeforeEach(func() {
			secondReturn = awshelper.GetDynamoDbClient()
		})

		It("should return the same client", func() {
			Expect(secondReturn).To(BeIdenticalTo(dynamoDbClient))
		})
	})
})
