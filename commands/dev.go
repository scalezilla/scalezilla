package commands

import (
	"context"
	"errors"

	"github.com/scalezilla/scalezilla/cluster"
	"github.com/scalezilla/scalezilla/logger"

	"github.com/urfave/cli/v3"
)

func Dev() *cli.Command {
	var (
		app  cluster.ClusterInitialConfig
		fail bool
	)

	return &cli.Command{
		Name:  "dev",
		Usage: "start a developer instance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host-ip-address",
				Aliases:     []string{"hip"},
				Value:       "127.0.0.1",
				Usage:       "IP address used by this instance",
				Destination: &app.HostIPAddress,
			},
			&cli.StringFlag{
				Name:        "bind-address",
				Aliases:     []string{"b"},
				Value:       "127.0.0.1",
				Usage:       "Bind address used by this instance",
				Destination: &app.BindAddress,
			},
			&cli.Uint16Flag{
				Name:        "http-port",
				Aliases:     []string{"hp"},
				Usage:       "http port to use",
				Value:       15000,
				Destination: &app.HTTPPort,
			},
			&cli.Uint16Flag{
				Name:        "raft-grpc-port",
				Aliases:     []string{"rgp"},
				Usage:       "grpc port for raft purpose",
				Value:       15001,
				Destination: &app.RaftGRPCPort,
			},
			&cli.Uint16Flag{
				Name:        "grpc-port",
				Aliases:     []string{"gp"},
				Usage:       "grpc port for internal purpose",
				Value:       15002,
				Destination: &app.GRPCPort,
			},
			&cli.BoolFlag{
				Name:        "fail",
				Hidden:      true,
				Destination: &fail,
			},
		},
		Action: func(context.Context, *cli.Command) error {
			if fail {
				return errors.New("test failure")
			}

			app.Logger = logger.NewLogger()
			cluster := cluster.NewCluster(app)

			return cluster.Start()
		},
	}
}
