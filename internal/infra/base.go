package infra

import (
	"github.com/dannydd88/dd-go"
	"github.com/urfave/cli/v2"
)

type GitUpContext struct {
	logger           dd.LevelLogger
	config           *Config
	workerPoolRunner *dd.WorkerPoolRunner
}

var globalContext GitUpContext

func AppInit(ctx *cli.Context) error {
	// ). init context
	globalContext = GitUpContext{}

	// ). init logger
	debug := ctx.Bool("debug")
	logLevel := dd.INFO
	if debug {
		logLevel = dd.DEBUG
	}
	globalContext.logger = dd.NewLevelLogger(logLevel)

	globalContext.logger.Debug("[app]", "AppInit finish")
	return nil
}

func CommandInit(ctx *cli.Context) error {
	// ). init config
	{
		globalContext.config = new(Config)
		err := dd.NewYAMLLoader[Config](dd.Ptr(ctx.String("config"))).Load(globalContext.config)
		if err != nil {
			return err
		}
	}

	// ). init ini config
	{
		iniPath := ctx.String("ini-config")
		profile := ctx.String("profile")
		repoConfig := globalContext.config.RepoConfig
		if repoConfig.Host == nil || repoConfig.Token == nil {
			iniConfig := new(INIConfig)
			err := dd.NewINILoader[INIConfig](dd.Ptr(iniPath), dd.Ptr(profile)).Load(iniConfig)
			if err == nil {
				if repoConfig.Host == nil {
					repoConfig.Host = iniConfig.Host
				}
				if repoConfig.Token == nil {
					repoConfig.Token = iniConfig.Token
				}
			} else {
				return err
			}
		}
	}

	// ). init pool runner
	globalContext.workerPoolRunner =
		dd.NewWorkerPoolRunner(&dd.WorkerPoolRunnerOptions{
			Logger: globalContext.logger,
		})

	globalContext.logger.Debug("[app]", "CommandInit finish")
	return nil
}

func GetLogger() dd.LevelLogger {
	return globalContext.logger
}

func GetConfig() *Config {
	return globalContext.config
}

func GetWorkerPoolRunner() dd.TaskRunner {
	return globalContext.workerPoolRunner
}
