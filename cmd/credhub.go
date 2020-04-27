package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	gviz "github.com/awalterschulze/gographviz"
	"github.com/urfave/cli/v2"
)

func CredHubCmd() *cli.Command {
	return &cli.Command{
		Name:    "credhub",
		Aliases: []string{"ch"},
		Subcommands: []*cli.Command{
			{
				Name:    "visualize-certificates",
				Aliases: []string{"vc"},
				Usage:   "Reads CredHub certificates from STDIN, writes Graphviz DOT to stdout",

				Action: func(c *cli.Context) error {
					graphStr, err := VisualizeCredHub(os.Stdin)

					if err != nil {
						return err
					}

					fmt.Println(graphStr)
					return nil
				},
			},
		},
	}
}

type credhubCertificate struct {
	Name     string  `json:"name"`
	SignedBy *string `json:"signed_by,omitempty"`
}

type credhubCertificates struct {
	Certificates []credhubCertificate `json:"certificates"`
}

func VisualizeCredHub(in io.Reader) (string, error) {
	var err error
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return "", err
	}

	var certificates credhubCertificates
	err = json.Unmarshal(b, &certificates)
	if err != nil {
		return "", err
	}

	g := gviz.NewEscape()

	err = g.SetName("G")
	if err != nil {
		return "", err
	}

	err = g.SetDir(true) // digraph not graph
	if err != nil {
		return "", err
	}

	for _, certificate := range certificates.Certificates {
		err = g.AddNode(
			"G",
			certificate.Name,
			map[string]string{"label": certificate.Name},
		)

		if err != nil {
			return "", err
		}
	}

	for _, certificate := range certificates.Certificates {
		if certificate.SignedBy == nil {
			// Do not produce an edge for certificates without parents
			continue
		}

		if *certificate.SignedBy == certificate.Name {
			// Do not produce an edge for self-signed certificates
			continue
		}

		err = g.AddEdge(
			certificate.Name,
			*certificate.SignedBy,
			/* directed */ true,
			/* attrs */ nil,
		)
		if err != nil {
			return "", err
		}
	}

	return g.String(), nil
}
