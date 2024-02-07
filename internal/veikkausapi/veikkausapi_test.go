package veikkausapi

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInternalVeikkausApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/veikkausapi-suite")
}
