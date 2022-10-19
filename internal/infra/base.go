package infra

import (
	"github.com/dannydd88/dd-go"
	"github.com/urfave/cli/v2"
)

type GitUpContext struct {
	logger           dd.Logger
	config           *Config
	workerPoolRunner *dd.WorkerPoolRunner
}

var globalContext GitUpContext

func Init(ctx *cli.Context) error {
	// ). init context
	globalContext = GitUpContext{}

	// ). init logger
	globalContext.logger = dd.NewDefaultLogger()

	// ). init config
	{
		c, err := loadConfig(dd.Ptr(ctx.String("config")))
		if err != nil {
			return err
		}
		globalContext.config = c
	}

	// ). init pool runner
	globalContext.workerPoolRunner =
		dd.NewWorkerPoolRunner(&dd.WorkerPoolRunnerOptions{
			Logger: globalContext.logger,
		})

	return nil
}

func GetLogger() dd.Logger {
	return globalContext.logger
}

func GetConfig() *Config {
	return globalContext.config
}

func GetWorkerPoolRunner() dd.TaskRunner {
	return globalContext.workerPoolRunner
}
