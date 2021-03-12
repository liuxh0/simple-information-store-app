package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationSuite")
}

var samCmd *exec.Cmd

var _ = BeforeSuite(func() {
	By("Starting local server")
	samCmd = exec.Command("sam", "local", "start-api")
	samCmd.Dir = ".."
	err := samCmd.Start()
	Expect(err).ShouldNot(HaveOccurred())

	By("Waiting until local server is ready")
	Eventually(func() error {
		_, err := http.Get("http://localhost:3000/hello")
		return err
	}, 10*time.Second).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	var err error

	err = samCmd.Process.Signal(os.Interrupt)
	Expect(err).ShouldNot(HaveOccurred())

	err = samCmd.Wait()
	Expect(err).ShouldNot(HaveOccurred())
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
		url := url.URL{
			Scheme: "http",
			Host:   "localhost:3000",
			Path:   "/i",
		}
		var err error
		resp, err = http.Post(url.String(), "", strings.NewReader(reqBody))
		Expect(err).ShouldNot(HaveOccurred())

		respBody = ReadReadCloserOrDie(resp.Body)
	})

	AfterEach(func() {
		if id, ok := getStringFromJsonString(respBody, "id"); ok && resp.StatusCode == 201 {
			fmt.Printf("Created with id %s\n", id)
			// TODO: Delete created value
		} else {
			fmt.Println("Nothing created")
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

func ReadReadCloserOrDie(rc io.ReadCloser) string {
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
