package main

import (
	"encoding/json"
	"errors"
	"simple-information-store-app/internal/service"
	"simple-information-store-app/internal/servicefakes"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("create-value handler", func() {
	const requestBody = "Test value in request body"

	var (
		fakeInfoCreator servicefakes.FakeInfoCreator
		handlerResponse events.APIGatewayProxyResponse
		handlerError    error
	)

	BeforeEach(func() {
		fakeInfoCreator = servicefakes.FakeInfoCreator{}
		infoCreator = &fakeInfoCreator
	})

	JustBeforeEach(func() {
		handlerResponse, handlerError = handler(events.APIGatewayProxyRequest{
			Body: requestBody,
		})
	})

	When("CreateInfo() returns no error", func() {
		BeforeEach(func() {
			fakeInfoCreator.CreateInfoCalls(func(id, value string) (service.Info, error) {
				return service.Info{
					ID:    id,
					Value: value,
				}, nil
			})
		})

		It("should work", func() {
			Expect(fakeInfoCreator.CreateInfoCallCount()).To(Equal(1))

			id, value := fakeInfoCreator.CreateInfoArgsForCall(0)
			Expect(id).To(HaveLen(36)) // A UUID should have 36 chars.
			Expect(value).To(Equal(requestBody))

			Expect(handlerResponse.StatusCode).To(Equal(201))

			responseBody := make(map[string]interface{})
			json.Unmarshal([]byte(handlerResponse.Body), &responseBody)
			Expect(responseBody["id"]).To(Equal(id))

			Expect(handlerError).To(BeNil())
		})

		It("should generate a new UUID each time", func() {
			handler(events.APIGatewayProxyRequest{Body: requestBody})

			Expect(fakeInfoCreator.CreateInfoCallCount()).To(Equal(2))
			id1, _ := fakeInfoCreator.CreateInfoArgsForCall(0)
			id2, _ := fakeInfoCreator.CreateInfoArgsForCall(1)

			Expect(id1).ToNot(Equal(id2))
		})
	})

	When("CreateInfo() returns ValueTooLongError", func() {
		BeforeEach(func() {
			fakeInfoCreator.CreateInfoCalls(func(id, value string) (service.Info, error) {
				return service.Info{}, service.ValueTooLongError{}
			})
		})

		It("should return 400 with error message", func() {
			Expect(handlerResponse.StatusCode).To(Equal(400))
			Expect(handlerResponse.Body).To(Equal(service.ValueTooLongError{}.Error()))
		})
	})

	When("CreateInfo() returns an error", func() {
		BeforeEach(func() {
			fakeInfoCreator.CreateInfoCalls(func(id, value string) (service.Info, error) {
				return service.Info{}, errors.New("")
			})
		})

		It("should return 500", func() {
			Expect(handlerResponse.StatusCode).To(Equal(500))
			Expect(handlerResponse.Body).To(BeEmpty())
		})
	})
})
