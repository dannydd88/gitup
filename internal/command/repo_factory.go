package command

import (
	"fmt"
	"strings"

	"gitup/internal/gitlab"
	"gitup/internal/infra"
	"gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
)

func buildRepoListor(config *infra.RepoConfig) (gitup.RepoListor, error) {
	var instance gitup.RepoListor
	var e error
	switch strings.ToLower(dd.Val(config.Type)) {
	case "gitlab":
		instance, e = gitlab.NewListor(config)
	default:
		return nil, fmt.Errorf("unsupport repostory type")
	}
	return instance, e
}

func buildRepoForker(config *infra.RepoConfig) (gitup.RepoForker, error) {
	var instance gitup.RepoForker
	var e error
	switch strings.ToLower(dd.Val(config.Type)) {
	case "gitlab":
		instance, e = gitlab.NewForker(config)
	default:
		return nil, fmt.Errorf("unsupport repostory type")
	}
	return instance, e
}
