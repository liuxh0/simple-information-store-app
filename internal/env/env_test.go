package env_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"simple-information-store-app/internal/env"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RunningInSamLocal()", func() {
	var ret bool

	BeforeEach(func() {
		UnsetEnvVars()
	})

	JustBeforeEach(func() {
		ret = env.RunningInSamLocal()
	})

	When("AWS_SAM_LOCAL environment variable is not set", func() {
		It("should return false", func() {
			Expect(ret).To(BeFalse())
		})
	})

	When("AWS_SAM_LOCAL environment variable is set", func() {
		BeforeEach(func() {
			setAwsSamLocalEnvVar()
		})

		It("should return true", func() {
			Expect(ret).To(BeTrue())
		})
	})
})

var _ = Describe("GetDynamoDbEndpoint()", func() {
	var ret string

	BeforeEach(func() {
		UnsetEnvVars()
	})

	JustBeforeEach(func() {
		ret = env.GetDynamoDbEndpoint()
	})

	When("AWS_SAM_LOCAL environment variable is set", func() {
		BeforeEach(func() {
			setAwsSamLocalEnvVar()
		})

		It("should return a local endpoint", func() {
			Expect(ret).To(Equal("http://host.docker.internal:8000"))
		})
	})

	When("GINKGO_TEST environment variable is set", func() {
		BeforeEach(func() {
			setGinkgoTestEnvVar()
		})

		It("should return a local endpoint", func() {
			Expect(ret).To(Equal("http://localhost:8000"))
		})
	})

	When("Both AWS_SAM_LOCAL and GINKGO_TEST are set", func() {
		BeforeEach(func() {
			setAwsSamLocalEnvVar()
			setGinkgoTestEnvVar()
		})

		Specify("AWS_SAM_LOCAL should take effect", func() {
			Expect(ret).To(Equal("http://host.docker.internal:8000"))
		})
	})

	When("Neither AWS_SAM_LOCAL nor GINKGO_TEST is set", func() {
		It("should return a empty string", func() {
			Expect(ret).To(BeEmpty())
		})
	})
})

var _ = Describe("GetValueTableName()", func() {
	const valueTableName = "test-ValueTable"

	var ret string

	BeforeEach(func() {
		err := os.Setenv("VALUE_TABLE_REF", valueTableName)
		Expect(err).ShouldNot(HaveOccurred())
		UnsetEnvVars()
	})

	JustBeforeEach(func() {
		ret = env.GetValueTableName()
	})

	When("AWS_SAM_LOCAL environment variable is set", func() {
		BeforeEach(func() {
			setAwsSamLocalEnvVar()
		})

		It("should return the table name for local DynamoDB", func() {
			Expect(ret).To(Equal(getLocalValueTableName()))
		})
	})

	When("GINKGO_TEST environment variable is set", func() {
		BeforeEach(func() {
			setGinkgoTestEnvVar()
		})

		It("should return the table name for local DynamoDB", func() {
			Expect(ret).To(Equal(getLocalValueTableName()))
		})
	})

	When("Neither AWS_SAM_LOCAL nor GINKGO_TEST is set", func() {
		It("should return the value of environment variable VALUE_TABLE_REF", func() {
			Expect(ret).To(Equal(valueTableName))
		})
	})
})

func setAwsSamLocalEnvVar() {
	err := os.Setenv("AWS_SAM_LOCAL", "true")
	if err != nil {
		panic(err)
	}
}

func setGinkgoTestEnvVar() {
	err := os.Setenv("GINKGO_TEST", "true")
	if err != nil {
		panic(err)
	}
}

func UnsetEnvVars() {
	var err error

	err = os.Unsetenv("AWS_SAM_LOCAL")
	if err != nil {
		panic(err)
	}

	err = os.Unsetenv("GINKGO_TEST")
	if err != nil {
		panic(err)
	}
}

func getLocalValueTableName() string {
	var err error

	bytes, err := ioutil.ReadFile("../../local-dynamodb-value-table.json")
	if err != nil {
		panic(err)
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		panic(err)
	}

	localTableName, ok := jsonMap["TableName"]
	if !ok {
		panic("TableName field not found")
	}

	return localTableName.(string)
}
