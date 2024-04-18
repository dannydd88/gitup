package command

import (
	"fmt"
	"strings"

	"github.com/dannydd88/gitup/internal/infra"
	"github.com/dannydd88/gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
)

func buildRepoList(config *infra.RepoConfig) (gitup.RepoList, error) {
	var instance gitup.RepoList
	var e error
	switch strings.ToLower(dd.Val(config.Type)) {
	case "gitlab":
		instance, e = gitup.NewGitlabList(&gitup.GitlabConfig{
			Host:           config.Host,
			Token:          config.Token,
			FilterArchived: config.FilterArchived,
		})
	default:
		return nil, fmt.Errorf("unsupport repostory type")
	}
	return instance, e
}

func buildRepoFork(config *infra.RepoConfig) (gitup.RepoFork, error) {
	var instance gitup.RepoFork
	var e error
	switch strings.ToLower(dd.Val(config.Type)) {
	case "gitlab":
		instance, e = gitup.NewGitlabFork(&gitup.GitlabConfig{
			Host:           config.Host,
			Token:          config.Token,
			FilterArchived: config.FilterArchived,
		})
	default:
		return nil, fmt.Errorf("unsupport repostory type")
	}
	return instance, e
}
