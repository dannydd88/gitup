package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dannydd88/gitup/internal/command"
	"github.com/dannydd88/gitup/internal/infra"

	"github.com/urfave/cli/v2"
)

var (
	version string = "dev"
	date    string = "dev"
	commit  string = "dev"
)

func main() {
	app := &cli.App{
		Name:    "gitup",
		Version: fmt.Sprintf("%s-%s (Build at %s)", version, commit, date),
		Usage:   "Tools for git repos management",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "gitup.yaml",
				Usage:   "Load config from yaml file",
			},
		},
		Before: infra.AppInit,
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
