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

func Init(ctx *cli.Context) error {
	// ). init context
	globalContext = GitUpContext{}

	// ). init logger
	globalContext.logger = dd.NewLevelLogger(dd.INFO)

	// ). init config
	{
		globalContext.config = &Config{}
		err := dd.NewYAMLLoader[Config](dd.Ptr(ctx.String("config"))).Load(globalContext.config)
		if err != nil {
			return err
		}
	}

	// ). init pool runner
	globalContext.workerPoolRunner =
		dd.NewWorkerPoolRunner(&dd.WorkerPoolRunnerOptions{
			Logger: globalContext.logger,
		})

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
