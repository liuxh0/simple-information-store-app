package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCreateValue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateValue Suite")
}
