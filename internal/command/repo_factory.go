package command

import (
	"fmt"
	"strings"

	"github.com/dannydd88/gitup/internal/gitlab"
	"github.com/dannydd88/gitup/internal/infra"
	"github.com/dannydd88/gitup/pkg/gitup"

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
