package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUpdateValue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateValue Suite")
}
