package gitlab

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
