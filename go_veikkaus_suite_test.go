package go_veikkaus_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoVeikkaus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoVeikkaus Suite")
}
