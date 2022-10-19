package main

import (
	"log"
	"os"

	"gitup/internal/command"
	"gitup/internal/infra"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gitup",
		Usage: "Tools for git repos management",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "gitup.yaml",
				Usage:   "Load config from yaml file",
			},
		},
		Before: infra.Init,
		Commands: []*cli.Command{
			command.NewSyncCommand(),
			command.NewForkCommand(),
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
