package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitup/internal/config"
	"gitup/internal/gitlab"
	"gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
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
				Value:   "config.json",
				Usage:   "Load config from json file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "sync",
				Usage: "Sync repo via config",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "concurrency",
						Aliases: []string{"c"},
						Usage:   "Git operation concurrency",
						Value:   10,
					},
				},
				Action: func(c *cli.Context) error {
					// ). Get current dir
					dir, err := os.Getwd()
					if err != nil {
						return err
					}

					// ). Load config
					configPath := c.String("config")
					if !filepath.IsAbs(configPath) {
						configPath = filepath.Join(dir, configPath)
					}
					config, err := config.LoadConfig(dd.String(configPath))
					if err != nil {
						return err
					}

					// ). Decide repository
					var r gitup.RepoHub
					switch strings.ToLower(config.RepoConfig.Type) {
					case "gitlab":
						r = gitlab.NewGitlab(&config.RepoConfig)
					default:
						return fmt.Errorf("[Main]Unsupport repostory type")
					}

					(&gitup.Runner{
						Hub:         r,
						Git:         config.GitConfig,
						Cwd:         config.Cwd,
						Concurrency: c.Int("concurrency"),
						Logger:      dd.NewDefaultLogger(),
					}).Execute()

					return nil
				},
			},
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
