package awshelper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAwshelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Awshelper Suite")
}
