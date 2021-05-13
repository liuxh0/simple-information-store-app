package integration_test

import (
	"fmt"
	"net/http"
	"simple-information-store-app/internal/service"
	"strings"

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
			service.NewInfoService().DeleteInfo(id)
		}
	})

	It("should return 201 with an id", func() {
		Expect(resp.StatusCode).To(Equal(201))

		bodyJsonMap, err := convertJsonStringToMap(respBody)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(bodyJsonMap).To(HaveKey("id"))
		Expect(bodyJsonMap["id"]).ShouldNot(BeEmpty())
	})

	When("request body is empty", func() {
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

		It("should contain an error message", func() {
			Expect(respBody).To(SatisfyAll(
				ContainSubstring("length"),
				ContainSubstring("1001"),
				ContainSubstring("max"),
				ContainSubstring("1000"),
			))
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
			id = generateNonExistingId()
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
			err := service.NewInfoService().DeleteInfo(id)
			if err != nil {
				panic(err)
			}

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
			id = generateNonExistingId()
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

			It("should contain an error message", func() {
				Expect(respBody).To(SatisfyAll(
					ContainSubstring("length"),
					ContainSubstring("1001"),
					ContainSubstring("max"),
					ContainSubstring("1000"),
				))
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
			err := service.NewInfoService().DeleteInfo(id)
			if err != nil {
				panic(err)
			}
		})

		It("should return 200", func() {
			Expect(resp.StatusCode).To(Equal(200))
			Expect(respBody).To(BeEmpty())

			By("checking if the value is updated", func() {
				info, err := service.NewInfoService().GetInfo(id)
				if err != nil {
					panic(err)
				}

				Expect(info.Value).To(Equal(reqBody))
			})
		})

		When("request body has more than 1000 characters", func() {
			BeforeEach(func() {
				reqBody = strings.Repeat("x", 1001)
			})

			It("should return 400", func() {
				Expect(resp.StatusCode).To(Equal(400))
			})

			It("should contain an error message", func() {
				Expect(respBody).To(SatisfyAll(
					ContainSubstring("length"),
					ContainSubstring("1001"),
					ContainSubstring("max"),
					ContainSubstring("1000"),
				))
			})
		})
	})
})

func generateNonExistingId() string {
	for {
		id := uuid.NewString()
		_, err := service.NewInfoService().GetInfo(id)
		switch err := err.(type) {
		case service.InfoNotFoundError:
			return id
		case nil:
			continue
		default:
			panic(err)
		}
	}
}
