package main

import (
	"errors"
	"simple-information-store-app/internal/service"
	"simple-information-store-app/internal/servicefakes"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get-value handler", func() {
	const (
		infoId    = "info-id"
		infoValue = "info value"
	)

	var (
		fakeInfoUpdater servicefakes.FakeInfoUpdater
		handlerResponse events.APIGatewayProxyResponse
	)

	BeforeEach(func() {
		fakeInfoUpdater = servicefakes.FakeInfoUpdater{}
		infoUpdater = &fakeInfoUpdater
	})

	JustBeforeEach(func() {
		var err error
		handlerResponse, err = handler(events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				"id": infoId,
			},
		})

		Expect(err).ShouldNot(HaveOccurred())
	})

	When("UpdateInfo() returns ValueTooLongError", func() {
		var valueTooLongError service.ValueTooLongError

		BeforeEach(func() {
			valueTooLongError = service.ValueTooLongError{}
			fakeInfoUpdater.UpdateInfoReturns(service.Info{}, valueTooLongError)
		})

		It("should return 400 with error message", func() {
			Expect(handlerResponse.StatusCode).To(Equal(400))

			errMsg := valueTooLongError.Error()
			Expect(errMsg).NotTo(BeEmpty())
			Expect(handlerResponse.Body).To(Equal(errMsg))
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})

	When("UpdateInfo() returns InfoNotFoundError", func() {
		BeforeEach(func() {
			fakeInfoUpdater.UpdateInfoReturns(service.Info{}, service.InfoNotFoundError{})
		})

		It("should return 404", func() {
			Expect(handlerResponse.StatusCode).To(Equal(404))
			Expect(handlerResponse.Body).To(BeEmpty())
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})

	When("UpdateInfo() returns an error", func() {
		BeforeEach(func() {
			fakeInfoUpdater.UpdateInfoReturns(service.Info{}, errors.New("error"))
		})

		It("should return 500", func() {
			Expect(handlerResponse.StatusCode).To(Equal(500))
			Expect(handlerResponse.Body).To(BeEmpty())
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})

	When("UpdateInfo() returns no error", func() {
		BeforeEach(func() {
			fakeInfoUpdater.UpdateInfoReturns(service.Info{}, nil)
		})

		It("should return 200", func() {
			Expect(handlerResponse.StatusCode).To(Equal(200))
			Expect(handlerResponse.Body).To(BeEmpty())
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})
})
