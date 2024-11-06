package templateutil

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTemplateutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Templateutil Suite")
}
