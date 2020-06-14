package timeline_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTimeline(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timeline Suite")
}
