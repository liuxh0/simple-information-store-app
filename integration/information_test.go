package integration_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
		endpointUrl := fmt.Sprintf("%s/i/%s", samHost, id)
		resp, err = http.Get(endpointUrl)
		Expect(err).ShouldNot(HaveOccurred())
		respBody = readReadCloserOrDie(resp.Body)
	})

	When("id does not exist", func() {
		BeforeEach(func() {
			for { // Find out an id that does not exist in the table
				id = uuid.NewString()
				result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
					TableName: &valueTableName,
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
			Expect(resp.StatusCode).To(Equal(201))
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

var _ = Describe("PUT /i/{id}", func() {
	var (
		id       string
		reqBody  string
		resp     *http.Response
		respBody string
	)

	BeforeEach(func() {
		id = ""
		reqBody = "An updated version of information updated by Integration test suite"
	})

	JustBeforeEach(func() {
		var err error

		Expect(id).ShouldNot(BeEmpty())
		endpointUrl := fmt.Sprintf("%s/i/%s", samHost, id)
		req, err := http.NewRequest(http.MethodPut, endpointUrl, strings.NewReader(reqBody))
		Expect(err).ShouldNot(HaveOccurred())

		httpClient := &http.Client{}
		resp, err = httpClient.Do(req)
		Expect(err).ShouldNot(HaveOccurred())
		respBody = readReadCloserOrDie(resp.Body)
	})

	When("id does not exist", func() {
		BeforeEach(func() {
			for { // Find out an id that does not exist in the table
				id = uuid.NewString()
				result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
					TableName: &valueTableName,
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

		When("request body has more than 1000 characters", func() {
			BeforeEach(func() {
				reqBody = strings.Repeat("x", 1001)
			})

			It("should return 400", func() {
				Expect(resp.StatusCode).To(Equal(400))
			})
		})
	})

	When("id exists", func() {
		const value = "First version created by Integration test suite"

		BeforeEach(func() { // Create a new item
			endpointUrl := fmt.Sprintf("%s/i", samHost)
			resp, err := http.Post(endpointUrl, "", strings.NewReader(value))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(201))
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

		It("should return 200", func() {
			Expect(resp.StatusCode).To(Equal(200))
			Expect(respBody).To(BeEmpty())

			By("checking if the value is updated", func() {
				result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
					TableName: &valueTableName,
					Key: map[string]*dynamodb.AttributeValue{
						"Id": {S: &id},
					},
				})

				if err != nil {
					panic(err)
				}

				Expect(*result.Item["Value"].S).To(Equal(reqBody))
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
})
