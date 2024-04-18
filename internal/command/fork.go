package command

import (
	"fmt"
	"os"

	"github.com/dannydd88/gitup/internal/infra"
	"github.com/dannydd88/gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func NewForkCommand() *cli.Command {
	return &cli.Command{
		Name:   "fork",
		Usage:  "Fork repo via config or flags",
		Before: infra.CommandInit,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "from-group",
				Aliases: []string{"fg"},
				Usage:   "Source repo's group",
			},
			&cli.StringFlag{
				Name:    "from-repo",
				Aliases: []string{"fr"},
				Usage:   "Source repo's name",
			},
			&cli.StringFlag{
				Name:    "to-group",
				Aliases: []string{"tg"},
				Usage:   "Target repo's group, can be null which means same group fork, should provide new repo name",
			},
			&cli.StringFlag{
				Name:    "to-repo",
				Aliases: []string{"tr"},
				Usage:   "Target repo's name, can be null which means using the same name as from-repo, but should not be same group fork",
			},
			&cli.StringFlag{
				Name:  "forks",
				Usage: "Fork config yaml file",
			},
			&cli.BoolFlag{
				Name:    "rm-fork-relation",
				Aliases: []string{"rfr"},
				Usage:   "Remove fork relationship",
			},
		},
		Action: func(ctx *cli.Context) error {
			config := infra.GetConfig()

			// ). check repo config
			if config == nil || config.RepoConfig == nil {
				return fmt.Errorf("[fork] missing repo config")
			}

			// ). decide repository type
			api, err := buildRepoFork(config.RepoConfig)
			if err != nil {
				return err
			}

			// ). prepare |ForkConfig| array
			var forkConfigs []*gitup.ForkConfig
			if existFlags(ctx, "forks") {
				// higher priority to use fork config file
				path := ctx.String("forks")
				if !dd.FileExists(dd.Ptr(path)) {
					return fmt.Errorf("[fork] cannot find config -> %s", path)
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				err = yaml.Unmarshal(data, &forkConfigs)
				if err != nil {
					return err
				}
			} else if existFlags(ctx, "from-group", "from-repo") {
				// ). prepare fromGroup and fromRepo
				fromGroup := ctx.String("from-group")
				fromRepo := ctx.String("from-repo")

				// ). build |ForkConfig|
				config := &gitup.ForkConfig{
					FromGroup: dd.Ptr(fromGroup),
					FromRepos: []*string{
						dd.Ptr(fromRepo),
					},
					RmForkRelation: dd.Ptr(ctx.Bool("rm-fork-relation")),
				}

				// ). modify |ForkConfig| according flags
				if ctx.IsSet("to-group") {
					toGroup := ctx.String("to-group")
					config.ToGroup = dd.Ptr(toGroup)
				}
				if ctx.IsSet("to-repo") {
					toRepo := ctx.String("to-repo")
					config.ToRepos = []*string{dd.Ptr(toRepo)}
				}

				// ). check |ForkConfig|
				if config.ToGroup == nil && len(config.ToRepos) == 0 {
					return fmt.Errorf(
						"[fork] ERROR: shoud provide one of flag %s | %s",
						"--to-group",
						"--to-repo",
					)
				}

				// ). append config
				forkConfigs = append(forkConfigs, config)
			} else {
				return fmt.Errorf(
					"[fork] ERROR: should provide fork info using file flag[%s] or cli flag[%s]",
					"--forks",
					"--from-group & --from-repo & (one of --to-group|--to-repo)",
				)
			}

			// ). construct forker and run
			(&gitup.Fork{
				Api:         api,
				ForkConfigs: forkConfigs,
				TaskRunner:  infra.GetWorkerPoolRunner(),
				Logger:      infra.GetLogger(),
			}).Go()

			return nil
		},
	}
}
