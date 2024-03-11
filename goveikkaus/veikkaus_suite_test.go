package goveikkaus

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVeikkaus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-veikkaus suite")
}
