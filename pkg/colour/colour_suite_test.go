package colour_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestColour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Colour Suite")
}
