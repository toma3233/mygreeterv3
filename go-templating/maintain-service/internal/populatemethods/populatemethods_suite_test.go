package populatemethods

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPopulatemethods(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PopulateMethods Suite")
}
