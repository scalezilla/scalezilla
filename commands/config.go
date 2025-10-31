package commands

import (
	"context"

	"github.com/scalezilla/scalezilla/cluster"
	"github.com/scalezilla/scalezilla/logger"

	"github.com/urfave/cli/v3"
)

func Config() *cli.Command {
	var (
		app      cluster.ClusterInitialConfig
		validate bool
	)

	return &cli.Command{
		Name:  "config",
		Usage: "start an instance with config file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "config file to use",
				Required:    true,
				Destination: &app.ConfigFile,
			},
			&cli.BoolFlag{
				Name:        "validate",
				Usage:       "check if your config is valid",
				Destination: &validate,
			},
			&cli.StringFlag{
				Name:        "test-raft-metric-prefix",
				Usage:       "override raft metric prefix during unit testing",
				Hidden:      true,
				Destination: &app.TestRaftMetricPrefix,
			},
			&cli.BoolFlag{
				Name:        "test",
				Usage:       "override some settings during unit testing",
				Hidden:      true,
				Destination: &app.Test,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			sigCtx, stop := cluster.BuildSignal(ctx)
			defer stop()

			app.Context = sigCtx
			app.Logger = logger.NewLogger()
			cluster, err := cluster.NewCluster(app)
			if err != nil {
				return err
			}

			if validate {
				app.Logger.Info().Msg("Your configuration is valid")
				return nil
			}

			return cluster.Start()
		},
	}
}
