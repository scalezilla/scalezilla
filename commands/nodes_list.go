package commands

import (
	"context"

	"github.com/scalezilla/scalezilla/cluster"
	"github.com/scalezilla/scalezilla/logger"

	"github.com/urfave/cli/v3"
)

func Nodes() *cli.Command {
	return &cli.Command{
		Name:    "nodes",
		Usage:   "cluster nodes options",
		Aliases: []string{"node"},
		Commands: []*cli.Command{
			list(),
		},
	}
}

func list() *cli.Command {
	var app cluster.NodesListHTTPConfig

	return &cli.Command{
		Name:  "list",
		Usage: "list options",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Aliases:     []string{"a"},
				Value:       "http://127.0.0.1:15000",
				Usage:       "HTTP(s) address to communicate with the cluster",
				Destination: &app.HTTPAddress,
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("SCALEZILLA_HTTP_ADDRESS"),
				),
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Usage:       "Token to use to bootstrap the cluster",
				Destination: &app.Token,
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("SCALEZILLA_TOKEN"),
				),
			},
			&cli.StringFlag{
				Name:        "kind",
				Aliases:     []string{"k"},
				Usage:       "Kind node to list, server or client",
				Destination: &app.Kind,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "output format can only be table or json",
				Value:       "table",
				Destination: &app.OutputFormat,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			if err := outputFormat(app.OutputFormat); err != nil {
				return err
			}
			sigCtx, stop := cluster.BuildSignal(ctx)
			defer stop()

			app.Context = sigCtx
			app.Logger = logger.NewLogger()
			return cluster.APICallsNodesList(app)
		},
	}
}
