package command

import (
	"gitup/internal/infra"
	"gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
	"github.com/urfave/cli/v2"
)

func NewSyncCommand() *cli.Command {
	return &cli.Command{
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
		Before: infra.Init,
		Action: func(c *cli.Context) error {
			config := infra.GetConfig()

			// ). Decide repository
			r, err := buildRepoHub(config.RepoConfig)
			if err != nil {
				return err
			}

			// ). construct syncer
			(&gitup.Syncer{
				Hub:         r,
				SyncConfig:  config.SyncConfig,
				Cwd:         config.Cwd,
				Concurrency: c.Int("concurrency"),
				Logger:      dd.NewDefaultLogger(),
			}).Go()

			return nil
		},
	}
}
