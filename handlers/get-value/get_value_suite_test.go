package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGetValue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GetValue Suite")
}
