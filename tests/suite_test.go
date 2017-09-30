package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"testing"
)

func TestBooklit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Booklit Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetLevel(logrus.FatalLevel)
})
