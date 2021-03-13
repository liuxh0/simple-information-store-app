package env_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"simple-information-store-app/internal/env"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const AwsSamLocalEnvVar = "AWS_SAM_LOCAL"

var _ = Describe("RunningInSamLocal()", func() {
	var ret bool

	JustBeforeEach(func() {
		ret = env.RunningInSamLocal()
	})

	When("AWS_SAM_LOCAL environment variable is not set", func() {
		BeforeEach(func() {
			err := os.Unsetenv(AwsSamLocalEnvVar)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return false", func() {
			Expect(ret).To(BeFalse())
		})
	})

	When("AWS_SAM_LOCAL environment variable is set", func() {
		BeforeEach(func() {
			err := os.Setenv(AwsSamLocalEnvVar, "true")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return true", func() {
			Expect(ret).To(BeTrue())
		})
	})
})

var _ = Describe("GetDynamoDbEndpoint()", func() {
	var ret string

	JustBeforeEach(func() {
		ret = env.GetDynamoDbEndpoint()
	})

	When("AWS_SAM_LOCAL environment variable is not set", func() {
		BeforeEach(func() {
			err := os.Unsetenv(AwsSamLocalEnvVar)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return a empty string", func() {
			Expect(ret).To(BeEmpty())
		})
	})

	When("AWS_SAM_LOCAL environment variable is set", func() {
		BeforeEach(func() {
			err := os.Setenv(AwsSamLocalEnvVar, "true")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return a local endpoint", func() {
			Expect(ret).To(Equal("http://host.docker.internal:8000"))
		})
	})
})

var _ = Describe("GetValueTableName()", func() {
	const valueTableName = "test-ValueTable"

	var ret string

	BeforeEach(func() {
		err := os.Setenv("VALUE_TABLE_REF", valueTableName)
		Expect(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		ret = env.GetValueTableName()
	})

	When("AWS_SAM_LOCAL environment variable is not set", func() {
		BeforeEach(func() {
			err := os.Unsetenv(AwsSamLocalEnvVar)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return the value of environment variable VALUE_TABLE_REF", func() {
			Expect(ret).To(Equal(valueTableName))
		})
	})

	When("AWS_SAM_LOCAL environment variable is set", func() {
		BeforeEach(func() {
			err := os.Setenv(AwsSamLocalEnvVar, "true")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return the table name for local DynamoDB", func() {
			var err error

			bytes, err := ioutil.ReadFile("../../local-dynamodb-value-table.json")
			Expect(err).ShouldNot(HaveOccurred())

			jsonMap := make(map[string]interface{})
			err = json.Unmarshal(bytes, &jsonMap)
			Expect(err).ShouldNot(HaveOccurred())

			localTableName, ok := jsonMap["TableName"]
			Expect(ok).To(BeTrue())

			Expect(ret).To(Equal(localTableName))
		})
	})
})
