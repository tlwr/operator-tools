package cmd_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tlwr/operator-tools/cmd"
)

var _ = Describe("YAML", func() {
	Describe("Traverse", func() {
		It("Should traverse YAML", func() {
			input := `
a:
  b:
    - name: first
      value: e
    - name: second
      value: f`

			output, err := cmd.TraverseYAML(
				"/a/b/name=first/value",
				strings.NewReader(input),
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal([]byte("e\n")))
		})

		It("Should throw an error for invalid YAML", func() {
			input := `
	this line starts with a tab and is therefore invalid`

			_, err := cmd.TraverseYAML("/path", strings.NewReader(input))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("line 2")))
			Expect(err).To(MatchError(ContainSubstring("found character that cannot start")))
		})

		It("Should throw an error for invalid path", func() {
			input := "this input does not matter"

			_, err := cmd.TraverseYAML("this is not a valid path", strings.NewReader(input))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("Invalid path")))
			Expect(err).To(MatchError(ContainSubstring("Expected to start with '/'")))
		})
	})
})
