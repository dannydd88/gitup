package command

import (
	"fmt"

	"github.com/dannydd88/gitup/internal/infra"
	"github.com/dannydd88/gitup/pkg/gitup"

	"github.com/urfave/cli/v2"
)

func NewSyncCommand() *cli.Command {
	return &cli.Command{
		Name:   "sync",
		Usage:  "Sync repo via config",
		Before: infra.CommandInit,
		Action: func(c *cli.Context) error {
			config := infra.GetConfig()

			// ). check config
			if config == nil || config.RepoConfig == nil || config.SyncConfig == nil {
				return fmt.Errorf("[Sync] gitup config error")
			}

			// ). decide repository type
			listor, err := buildRepoListor(config.RepoConfig)
			if err != nil {
				return err
			}

			// ). construct syncer and run
			syncConfig := &gitup.SyncConfig{
				Bare:   config.SyncConfig.Bare,
				Groups: config.SyncConfig.Groups,
			}
			(&gitup.Syncer{
				Api:        listor,
				SyncConfig: syncConfig,
				Cwd:        config.Cwd,
				TaskRunner: infra.GetWorkerPoolRunner(),
				Logger:     infra.GetLogger(),
			}).Go()

			return nil
		},
	}
}
