package gitup

import (
	"fmt"

	"github.com/dannydd88/dd-go"
	gitlabapi "github.com/xanzy/go-gitlab"
)

const (
	baseURL = "https://%s/api/v4"
)

type gitlabContext struct {
	apiClient *gitlabapi.Client
}

func newGitlabClient(token, host *string) (*gitlabapi.Client, error) {
	c, err := gitlabapi.NewClient(
		dd.Val(token),
		gitlabapi.WithBaseURL(fmt.Sprintf(baseURL, dd.Val(host))),
	)
	return c, err
}

type GitlabConfig struct {
	Host           *string
	Token          *string
	FilterArchived bool
}

// NewGitlabList
// Helper function to create |RepoList| gitlab implement
func NewGitlabList(config *GitlabConfig) (RepoList, error) {
	// ). construct gitlab client
	c, err := newGitlabClient(config.Token, config.Host)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabList{
		gitlabContext: gitlabContext{
			apiClient: c,
		},
		projects:       make(map[string][]*Repo),
		filterArchived: config.FilterArchived,
	}
	return g, nil
}

// NewGitlabFork
// Helper function to create |RepoFork| gitlab implement
func NewGitlabFork(config *GitlabConfig) (RepoFork, error) {
	// ). construct gitlab client
	c, err := newGitlabClient(config.Token, config.Host)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabFork{
		gitlabList: gitlabList{
			gitlabContext: gitlabContext{
				apiClient: c,
			},
			projects:       make(map[string][]*Repo),
			filterArchived: config.FilterArchived,
		},
	}

	return g, nil
}
