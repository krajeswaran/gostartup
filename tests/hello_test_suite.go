package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestHello(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hello Service Suite")
}
