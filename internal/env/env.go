package env

import "os"

// RunningInSamLocal returns if it is running in sam local environment.
func RunningInSamLocal() bool {
	_, ok := os.LookupEnv("AWS_SAM_LOCAL")
	return ok
}

func GetDynamoDbEndpoint() string {
	if RunningInSamLocal() {
		return "http://docker.for.mac.localhost:8000"
	}
	return ""
}

func GetValueTableName() string {
	if RunningInSamLocal() {
		return "simple-information-store-app-local-ValueTable"
	}
	return os.Getenv("VALUE_TABLE_REF")
}
