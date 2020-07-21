package main

import (
	"os"

	"github.com/tlwr/operator-tools/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.CredHubCmd(),
			cmd.HTTPCmd(),
			cmd.X509Cmd(),
			cmd.YamlCmd(),
		},
	}

	app.Run(os.Args)
}
