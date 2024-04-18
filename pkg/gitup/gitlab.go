package gitup

import (
	"fmt"

	"github.com/dannydd88/dd-go"
	gitlabapi "github.com/xanzy/go-gitlab"
)

const (
	baseURL = "https://%s/api/v4"
)

type GitlabApi interface {
	Api() *gitlabapi.Client
}

func NewGitlabApi(token, host *string) (GitlabApi, error) {
	// ). construct gitlab client
	c, err := gitlabapi.NewClient(
		dd.Val(token),
		gitlabapi.WithBaseURL(fmt.Sprintf(baseURL, dd.Val(host))),
	)
	if err != nil {
		return nil, err
	}

	api := &gitlabContext{
		apiClient: c,
	}

	return api, nil
}

type gitlabContext struct {
	apiClient *gitlabapi.Client
}

func (g *gitlabContext) Api() *gitlabapi.Client {
	return g.apiClient
}

type GitlabConfig struct {
	Host           *string
	Token          *string
	FilterArchived bool
}

// NewGitlabList
// Helper function to create |RepoList| gitlab implement
func NewGitlabList(config *GitlabConfig) (RepoList, error) {
	// ). construct |GitlabApi|
	api, err := NewGitlabApi(config.Token, config.Host)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabList{
		gitlab:         api,
		projects:       make(map[string][]*Repo),
		filterArchived: config.FilterArchived,
	}
	return g, nil
}

// NewGitlabFork
// Helper function to create |RepoFork| gitlab implement
func NewGitlabFork(config *GitlabConfig) (RepoFork, error) {
	// ). construct |GitlabApi|
	api, err := NewGitlabApi(config.Token, config.Host)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabFork{
		gitlabList: gitlabList{
			gitlab:         api,
			projects:       make(map[string][]*Repo),
			filterArchived: config.FilterArchived,
		},
	}

	return g, nil
}
