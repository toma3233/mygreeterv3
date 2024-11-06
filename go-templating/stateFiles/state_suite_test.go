package stateFiles

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStateFiles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateFiles Suite")
}
