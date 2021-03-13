package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"simple-information-store-app/internal/env"
	"simple-information-store-app/internal/helper/awshelper"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
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
})

var _ = Describe("POST /i", func() {
	var (
		reqBody  string
		resp     *http.Response
		respBody string
	)

	BeforeEach(func() {
		reqBody = "A piece of information created by Integration test suite"
	})

	JustBeforeEach(func() {
		var err error

		endpointUrl := fmt.Sprintf("%s/i", samHost)
		resp, err = http.Post(endpointUrl, "", strings.NewReader(reqBody))
		Expect(err).ShouldNot(HaveOccurred())

		respBody = readReadCloserOrDie(resp.Body)
		if id, ok := getStringFromJsonString(respBody, "id"); ok && resp.StatusCode == 201 {
			fmt.Printf("Created item with id %s\n", id)
		}
	})

	AfterEach(func() {
		if id, ok := getStringFromJsonString(respBody, "id"); ok && resp.StatusCode == 201 {
			deleteItem(id)
		}
	})

	It("should return 201 with an id", func() {
		Expect(resp.StatusCode).To(Equal(201))

		bodyJsonMap, err := convertJsonStringToMap(respBody)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(bodyJsonMap).To(HaveKey("id"))
		Expect(bodyJsonMap["id"]).ShouldNot(BeEmpty())
	})

	PWhen("request body is empty", func() {
		BeforeEach(func() {
			reqBody = ""
		})

		It("should return 201", func() {
			Expect(resp.StatusCode).To(Equal(201))
		})
	})

	When("request body has more than 1000 characters", func() {
		BeforeEach(func() {
			reqBody = strings.Repeat("x", 1001)
		})

		It("should return 400", func() {
			Expect(resp.StatusCode).To(Equal(400))
		})
	})
})

var _ = Describe("GET /i/{id}", func() {
	var (
		id       string
		resp     *http.Response
		respBody string
	)

	BeforeEach(func() {
		id = ""
	})

	JustBeforeEach(func() {
		var err error

		Expect(id).ShouldNot(BeEmpty())
		endpintUrl := fmt.Sprintf("%s/i/%s", samHost, id)
		resp, err = http.Get(endpintUrl)
		Expect(err).ShouldNot(HaveOccurred())
		respBody = readReadCloserOrDie(resp.Body)
	})

	When("id does not exist", func() {
		BeforeEach(func() {
			for { // Find out an id that does not exist in the table
				id = uuid.NewString()

				dynamoDbClient := awshelper.GetDynamoDbClient(dynamoDbEndpoint)
				tableName := env.GetValueTableName()
				result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
					TableName: &tableName,
					Key: map[string]*dynamodb.AttributeValue{
						"Id": {S: &id},
					},
				})
				Expect(err).ShouldNot(HaveOccurred())

				if result.Item == nil {
					break
				}
			}
		})

		It("should return 404", func() {
			Expect(resp.StatusCode).To(Equal(404))
			Expect(respBody).To(BeEmpty())
		})
	})

	When("id exists", func() {
		const value = "Test value used by Integration test suite"

		BeforeEach(func() { // Create a new item
			endpointUrl := fmt.Sprintf("%s/i", samHost)
			resp, err := http.Post(endpointUrl, "", strings.NewReader(value))
			Expect(err).ShouldNot(HaveOccurred())
			respBody := readReadCloserOrDie(resp.Body)
			newId, ok := getStringFromJsonString(respBody, "id")
			Expect(ok).To(BeTrue())
			fmt.Printf("Created item with id %s\n", newId)

			id = newId
		})

		AfterEach(func() { // Delete the new item created for the test
			err := deleteItem(id)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return 200 with the value", func() {
			Expect(resp.StatusCode).To(Equal(200))
			Expect(respBody).To(Equal(value))
		})
	})
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
	dynamoDbClient := awshelper.GetDynamoDbClient(dynamoDbEndpoint)
	tableName := env.GetValueTableName()
	_, err := dynamoDbClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: &id},
		},
	})
	return err
}
