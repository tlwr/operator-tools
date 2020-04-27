package cmd_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tlwr/operator-tools/cmd"
)

var _ = Describe("CredHub", func() {
	Describe("Visualize", func() {
		It("Should produce a digraph", func() {
			input := `{"certificates": [{
        "name": "/my/cert",
				"signed_by": "/my/ca"
			}, {
        "name": "/my/ca",
				"signed_by": "/my/ca"
			}, {
        "name": "/my/leaf"
			}]}`

			graph, err := cmd.VisualizeCredHub(strings.NewReader(input))
			Expect(err).NotTo(HaveOccurred())

			Expect(graph).To(Equal(`digraph G {
	"/my/ca"->"/my/cert";
	"/my/ca" [ label="/my/ca" ];
	"/my/cert" [ label="/my/cert" ];
	"/my/leaf" [ label="/my/leaf" ];

}
`))
		})
	})
})
