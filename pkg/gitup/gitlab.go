package gitup

import (
	"fmt"

	"github.com/dannydd88/dd-go"
	gitlabapi "github.com/xanzy/go-gitlab"
)

const (
	baseURL = "https://%s/api/v4"

	TagGitlab = "[gitlab]"
)

type GitlabApi interface {
	// Api - Returen the |gitlab| api instance for gitlab api access
	Api() *gitlabapi.Client

	// Logger - Return the current logger for logging
	Logger() dd.LevelLogger
}

func NewGitlabApi(token, host *string, logger dd.LevelLogger) (GitlabApi, error) {
	// ). construct gitlab client
	c, err := gitlabapi.NewClient(
		dd.Val(token),
		gitlabapi.WithBaseURL(fmt.Sprintf(baseURL, dd.Val(host))),
	)
	if err != nil {
		return nil, err
	}

	api := &gitlabContext{
		api:    c,
		logger: logger,
	}

	return api, nil
}

type gitlabContext struct {
	api    *gitlabapi.Client
	logger dd.LevelLogger
}

func (g *gitlabContext) Api() *gitlabapi.Client {
	return g.api
}

func (g *gitlabContext) Logger() dd.LevelLogger {
	return g.logger
}

type GitlabConfig struct {
	Host           *string
	Token          *string
	FilterArchived bool
	Logger         dd.LevelLogger
}

// NewGitlabList
// Helper function to create |RepoList| gitlab implement
func NewGitlabList(config *GitlabConfig) (RepoList, error) {
	// ). construct |GitlabApi|
	api, err := NewGitlabApi(config.Token, config.Host, config.Logger)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabList{
		GitlabApi:      api,
		projects:       make(map[string][]*Repo),
		filterArchived: config.FilterArchived,
	}
	return g, nil
}

// NewGitlabFork
// Helper function to create |RepoFork| gitlab implement
func NewGitlabFork(config *GitlabConfig) (RepoFork, error) {
	// ). construct |GitlabApi|
	api, err := NewGitlabApi(config.Token, config.Host, config.Logger)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabFork{
		gitlabList: gitlabList{
			GitlabApi:      api,
			projects:       make(map[string][]*Repo),
			filterArchived: config.FilterArchived,
		},
	}

	return g, nil
}
