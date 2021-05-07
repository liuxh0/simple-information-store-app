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
		fakeInfoGetter  servicefakes.FakeInfoGetter
		handlerResponse events.APIGatewayProxyResponse
	)

	BeforeEach(func() {
		fakeInfoGetter = servicefakes.FakeInfoGetter{}
		infoGetter = &fakeInfoGetter
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

	When("GetInfo() returns InfoNotFoundError", func() {
		BeforeEach(func() {
			fakeInfoGetter.GetInfoReturns(service.Info{}, service.InfoNotFoundError{})
		})

		It("should return 404", func() {
			Expect(handlerResponse.StatusCode).To(Equal(404))
			Expect(handlerResponse.Body).To(BeEmpty())
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})

	When("GetInfo() returns an error", func() {
		BeforeEach(func() {
			fakeInfoGetter.GetInfoReturns(service.Info{}, errors.New("error"))
		})

		It("should return 500", func() {
			Expect(handlerResponse.StatusCode).To(Equal(500))
			Expect(handlerResponse.Body).To(BeEmpty())
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})

	When("GetInfo() returns no error", func() {
		BeforeEach(func() {
			fakeInfoGetter.GetInfoReturns(service.Info{Value: infoValue}, nil)
		})

		It("should return 200 with body", func() {
			Expect(handlerResponse.StatusCode).To(Equal(200))
			Expect(handlerResponse.Body).To(Equal(infoValue))
			Expect(handlerResponse.Headers).To(BeEmpty())
		})
	})
})
