package main

import (
	"context"
	"os"

	"github.com/scalezilla/scalezilla/commands"
	"github.com/scalezilla/scalezilla/logger"

	"github.com/urfave/cli/v3"
)

func main() {
	usage := "A new cli manage your platforms and clusters"
	description := "Scalezilla is an orchestrator for managing containerized applications across multiple hosts without heaviness"

	cmd := cli.Command{
		Name:                  "scalezilla",
		Usage:                 usage,
		Description:           description,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			commands.Agent(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.NewLogger().Fatal().Err(err).Msg("Error occured while executing the program")
		// MUST keep os.Exit otherwise testing fatal won't work
		os.Exit(1)
	}
}
