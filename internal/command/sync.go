package command

import (
	"fmt"
	"gitup/internal/infra"
	"gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
	"github.com/urfave/cli/v2"
)

func NewSyncCommand() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "Sync repo via config",
		Action: func(c *cli.Context) error {
			config := infra.GetConfig()

			// ). check config
			if config == nil || config.RepoConfig == nil || config.SyncConfig == nil {
				return fmt.Errorf("[Main] gitup config error")
			}

			// ). decide repository type
			listor, err := buildRepoListor(config.RepoConfig)
			if err != nil {
				return err
			}

			// ). construct syncer and run
			(&gitup.Syncer{
				Api:        listor,
				SyncConfig: config.SyncConfig,
				Cwd:        config.Cwd,
				Logger:     dd.NewDefaultLogger(),
			}).Go()

			return nil
		},
	}
}
