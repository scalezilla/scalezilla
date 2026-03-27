package commands

import (
	"context"

	"github.com/scalezilla/scalezilla/cluster"
	"github.com/scalezilla/scalezilla/logger"

	"github.com/urfave/cli/v3"
)

func Deployment() *cli.Command {
	return &cli.Command{
		Name:    "deployment",
		Usage:   "deployment options",
		Aliases: []string{"deploy"},
		Commands: []*cli.Command{
			deploymentApply(),
		},
	}
}

func deploymentApply() *cli.Command {
	var app cluster.DeploymentApplyHTTPConfig

	return &cli.Command{
		Name:  "apply",
		Usage: "apply options",
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
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "deployment file to apply",
				Required:    true,
				Destination: &app.File,
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
			sigCtx, stop := cluster.BuildSignal(ctx)
			defer stop()

			app.Context = sigCtx
			app.Logger = logger.NewLogger()
			return cluster.APICallsDeploymentApply(app)
		},
	}
}
