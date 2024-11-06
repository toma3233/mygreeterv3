package resourcelinks

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResourceLinks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ResourceLinks Suite")
}
