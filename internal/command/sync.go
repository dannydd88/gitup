package command

import (
	"fmt"

	"github.com/dannydd88/dd-go"
	"github.com/dannydd88/gitup/internal/infra"
	"github.com/dannydd88/gitup/pkg/gitup"

	"github.com/urfave/cli/v2"
)

func NewSyncCommand() *cli.Command {
	return &cli.Command{
		Name:   "sync",
		Usage:  "Sync repo via config",
		Before: infra.CommandInit,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Usage:   "Groups that need to sync [higher priority than sync settings in yaml file]",
			},
			&cli.BoolFlag{
				Name:  "bare",
				Usage: "Should sync repo in bare way",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			config := infra.GetConfig()

			// ). check repo config
			if config == nil || config.RepoConfig == nil || config.SyncConfig == nil {
				return fmt.Errorf("[Sync] missing repo config")
			}

			// ). decide repository type
			api, err := buildRepoList(config.RepoConfig)
			if err != nil {
				return err
			}

			// ). check and build sync config
			var syncConfig *gitup.SyncConfig
			if existFlags(c, "group") {
				// higher priority to use cli flag
				syncConfig = &gitup.SyncConfig{
					Token:  config.RepoConfig.Token,
					Bare:   c.Bool("bare"),
					Groups: dd.PtrSlice(c.StringSlice("group")),
				}
			} else if config.SyncConfig != nil {
				syncConfig = &gitup.SyncConfig{
					Token:  config.RepoConfig.Token,
					Bare:   config.SyncConfig.Bare,
					Groups: config.SyncConfig.Groups,
				}
			} else {
				return fmt.Errorf(
					"[Sync] ERROR: should provide sync info using cli flag[%s] or in config file section[%s]",
					"--group",
					"sync",
				)
			}

			// ). construct syncer and run
			(&gitup.Sync{
				Api:        api,
				SyncConfig: syncConfig,
				Cwd:        config.Cwd,
				TaskRunner: infra.GetWorkerPoolRunner(),
				Logger:     infra.GetLogger(),
			}).Go()

			return nil
		},
	}
}
