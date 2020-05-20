package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cppforlife/go-patch/patch"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func YamlCmd() *cli.Command {
	return &cli.Command{
		Name:    "yaml",
		Aliases: []string{"y"},
		Subcommands: []*cli.Command{
			{
				Name:    "find",
				Aliases: []string{"f"},
				Usage:   "Traverses a YAML file using BOSH interpolation syntax",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Required: true,
					},
				},

				Action: func(c *cli.Context) error {
					path := c.String("path")
					if path == "" {
						cli.ShowAppHelpAndExit(c, 1)
					}

					outputBytes, err := TraverseYAML(path, os.Stdin)
					if err != nil {
						return err
					}

					os.Stdout.Write(outputBytes)

					return nil
				},
			},
		},
	}
}

func TraverseYAML(path string, r io.Reader) ([]byte, error) {
	pathPointer, err := patch.NewPointerFromString(path)
	if err != nil {
		return []byte{}, fmt.Errorf("Invalid path: %w", err)
	}

	inputBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}

	var input interface{}
	err = yaml.Unmarshal(inputBytes, &input)
	if err != nil {
		return []byte{}, err
	}

	output, err := patch.FindOp{Path: pathPointer}.Apply(input)
	if err != nil {
		return []byte{}, err
	}

	outputBytes, err := yaml.Marshal(output)
	if err != nil {
		return []byte{}, err
	}

	return outputBytes, nil
}
