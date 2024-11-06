package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPopulateMethods(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MaintainService Suite")
}
