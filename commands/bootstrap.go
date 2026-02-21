package commands

import (
	"context"

	"github.com/scalezilla/scalezilla/cluster"
	"github.com/scalezilla/scalezilla/logger"

	"github.com/urfave/cli/v3"
)

func Bootstrap() *cli.Command {
	return &cli.Command{
		Name:  "bootstrap",
		Usage: "cluster bootstrap options",
		Commands: []*cli.Command{
			status(),
			bootstrap(),
		},
	}
}

func status() *cli.Command {
	var app cluster.ClusterHTTPCallBaseConfig

	return &cli.Command{
		Name:  "status",
		Usage: "status options",
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
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			sigCtx, stop := cluster.BuildSignal(ctx)
			defer stop()

			app.Context = sigCtx
			app.Logger = logger.NewLogger()
			cluster.APICallsBootstrapStatus(app)

			return nil
		},
	}
}

func bootstrap() *cli.Command {
	var app cluster.BootstrapClusterHTTPConfig

	return &cli.Command{
		Name:  "cluster",
		Usage: "bootstrap options",
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
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			sigCtx, stop := cluster.BuildSignal(ctx)
			defer stop()

			app.Context = sigCtx
			app.Logger = logger.NewLogger()
			return cluster.APICallsBootstrapCluster(app)
		},
	}
}
