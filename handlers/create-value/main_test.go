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
	)

	BeforeEach(func() {
		fakeInfoCreator = servicefakes.FakeInfoCreator{}
		infoCreator = &fakeInfoCreator
	})

	JustBeforeEach(func() {
		var err error
		handlerResponse, err = handler(events.APIGatewayProxyRequest{
			Body: requestBody,
		})

		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should call CreateInfo() with a generated UUID", func() {
		Expect(fakeInfoCreator.CreateInfoCallCount()).To(Equal(1))

		id, value := fakeInfoCreator.CreateInfoArgsForCall(0)
		Expect(id).To(HaveLen(36)) // A UUID should have 36 chars.
		Expect(value).To(Equal(requestBody))
	})

	It("should generate a new UUID each time", func() {
		handler(events.APIGatewayProxyRequest{Body: requestBody})

		Expect(fakeInfoCreator.CreateInfoCallCount()).To(Equal(2))
		id1, _ := fakeInfoCreator.CreateInfoArgsForCall(0)
		id2, _ := fakeInfoCreator.CreateInfoArgsForCall(1)

		Expect(id1).ToNot(Equal(id2))
	})

	When("CreateInfo() returns ValueTooLongError", func() {
		var valueTooLongError service.ValueTooLongError

		BeforeEach(func() {
			valueTooLongError = service.ValueTooLongError{}
			fakeInfoCreator.CreateInfoReturns(service.Info{}, valueTooLongError)
		})

		It("should return 400 with error message", func() {
			Expect(handlerResponse.StatusCode).To(Equal(400))

			errMsg := valueTooLongError.Error()
			Expect(errMsg).NotTo(BeEmpty())
			Expect(handlerResponse.Body).To(Equal(errMsg))
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})

	When("CreateInfo() returns an error", func() {
		BeforeEach(func() {
			fakeInfoCreator.CreateInfoReturns(service.Info{}, errors.New("error"))
		})

		It("should return 500", func() {
			Expect(handlerResponse.StatusCode).To(Equal(500))
			Expect(handlerResponse.Body).To(BeEmpty())
			Expect(handlerResponse.Headers).To(BeEmpty())
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

		It("should return 201 with Id", func() {
			Expect(handlerResponse.StatusCode).To(Equal(201))
			Expect(handlerResponse.Headers).To(BeEmpty())

			responseBody := make(map[string]interface{})
			json.Unmarshal([]byte(handlerResponse.Body), &responseBody)

			id, _ := fakeInfoCreator.CreateInfoArgsForCall(0)
			Expect(responseBody["id"]).To(Equal(id))
		})
	})
})
