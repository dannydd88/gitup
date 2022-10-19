package main

import (
	"log"
	"os"

	"gitup/internal/command"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gitup",
		Usage: "git update according config file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.yaml",
				Usage:   "Load config from yaml file",
			},
		},
		Commands: []*cli.Command{
			command.NewSyncCommand(),
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
