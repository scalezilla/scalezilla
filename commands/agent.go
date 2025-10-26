package commands

import (
	"github.com/urfave/cli/v3"
)

func Agent() *cli.Command {
	return &cli.Command{
		Name:  "agent",
		Usage: "option to start an instance",
		Commands: []*cli.Command{
			Dev(),
		},
	}
}
