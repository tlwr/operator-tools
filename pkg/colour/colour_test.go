package colour_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/tlwr/operator-tools/pkg/colour"
)

var _ = Describe("Colour", func() {
	DescribeTable(
		"foreground colours",
		func(input string, expected string, fn func(string) string) {
			Expect(fn(input)).To(Equal(expected))
		},
		Entry("red", "red", "\033[31mred\033[0m", colour.Red),
		Entry("green", "green", "\033[32mgreen\033[0m", colour.Green),
		Entry("yellow", "yellow", "\033[33myellow\033[0m", colour.Yellow),
		Entry("blue", "blue", "\033[34mblue\033[0m", colour.Blue),
		Entry("magenta", "magenta", "\033[35mmagenta\033[0m", colour.Magenta),
		Entry("cyan", "cyan", "\033[36mcyan\033[0m", colour.Cyan),
	)
})
