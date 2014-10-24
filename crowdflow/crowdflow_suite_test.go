package crowdflow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCrowdflow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crowdflow Suite")
}
