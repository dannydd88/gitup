package command

import (
	"fmt"
	"gitup/internal/infra"
	"gitup/pkg/gitup"
	"os"

	"github.com/dannydd88/dd-go"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func NewForkCommand() *cli.Command {
	return &cli.Command{
		Name:  "fork",
		Usage: "Fork repo via config or flags",
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
				Usage:   "Target repo's group",
			},
			&cli.StringFlag{
				Name:    "to-repo",
				Aliases: []string{"tr"},
				Usage:   "Target repo's name",
			},
			&cli.StringFlag{
				Name:  "forks",
				Usage: "Fork config yaml file",
			},
		},
		Action: func(ctx *cli.Context) error {
			config := infra.GetConfig()

			// ). decide repository type
			forker, err := buildRepoForker(config.RepoConfig)
			if err != nil {
				return err
			}

			// ). prepare |ForkConfig| array
			var forkConfigs []*gitup.ForkConfig
			if ctx.IsSet("forks") {
				// higher priority to use fork config file
				path := ctx.String("forks")
				if !dd.FileExists(dd.Ptr(path)) {
					return fmt.Errorf("[Main] Cannot find config -> %s", path)
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				err = yaml.Unmarshal(data, &forkConfigs)
				if err != nil {
					return err
				}
			} else {
				// individual repo fork, check flags
				if !ctx.IsSet("from-group") || !ctx.IsSet("from-repo") || !ctx.IsSet("to-group") {
					return fmt.Errorf("[Main] Missing one of flag(from-group/from-repo/to-group)")
				}
				config := &gitup.ForkConfig{
					FromGroup: dd.Ptr(ctx.String("from-group")),
					FromRepos: []*string{
						dd.Ptr(ctx.String("from-repo")),
					},
					ToGroup: dd.Ptr(ctx.String("to-group")),
				}
				if ctx.IsSet("to-repo") {
					config.ToRepos = []*string{dd.Ptr(ctx.String("to-repo"))}
				}
				forkConfigs = append(forkConfigs, config)
			}

			// ). construct forker and run
			(&gitup.Forker{
				Api:         forker,
				ForkConfigs: forkConfigs,
				Logger:      dd.NewDefaultLogger(),
			}).Go()

			return nil
		},
	}
}
