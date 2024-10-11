package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
			&cli.StringFlag{
				Name:  "ini-config",
				Value: generateDefaultINIPath(),
				Usage: "Load profile config from ini file",
			},
			&cli.StringFlag{
				Name:  "profile",
				Value: "default",
				Usage: "Target profile that need to read from ini file",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "Debug mode",
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

func generateDefaultINIPath() string {
	iniConfig := "gitup.ini"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return iniConfig
	}

	return filepath.Join(homeDir, ".config", iniConfig)
}
