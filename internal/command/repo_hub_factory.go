package command

import (
	"fmt"
	"strings"

	"gitup/internal/gitlab"
	"gitup/internal/infra"
	"gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
)

func buildRepoHub(config *infra.RepoConfig) (gitup.RepoHub, error) {
	var r gitup.RepoHub
	switch strings.ToLower(dd.Val(config.Type)) {
	case "gitlab":
		r = gitlab.NewGitlab(config)
	default:
		return nil, fmt.Errorf("unsupport repostory type")
	}
	return r, nil
}
