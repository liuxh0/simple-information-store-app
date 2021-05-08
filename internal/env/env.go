package env

import (
	"os"
)

// RunningInSamLocal returns if it is running in sam local environment.
func RunningInSamLocal() bool {
	_, ok := os.LookupEnv("AWS_SAM_LOCAL")
	return ok
}

func runningInGinkgoTest() bool {
	_, ok := os.LookupEnv("GINKGO_TEST")
	return ok
}

// GetDynamoDbEndpoint returns the DynamoDB endpoint according to running environment.
func GetDynamoDbEndpoint() string {
	if RunningInSamLocal() {
		return "http://dynamodb:8000"
	} else if runningInGinkgoTest() {
		return "http://localhost:8000"
	}
	return ""
}

// GetValueTableName returns the name for ValueTable according to running environment.
func GetValueTableName() string {
	if RunningInSamLocal() || runningInGinkgoTest() {
		return "simple-information-store-app-local-ValueTable"
	}
	return os.Getenv("VALUE_TABLE_REF")
}
