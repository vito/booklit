package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBooklit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Booklit Suite")
}
