package gitlab

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitup/internal/infra"
	"gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
	gitlabapi "github.com/xanzy/go-gitlab"
)

type gitlab struct {
	host           *string
	token          *string
	projects       map[string][]*gitup.Repo
	filterArchived bool
}

// NewGitlab -
func NewGitlab(config *infra.RepoConfig) gitup.RepoHub {
	g := &gitlab{
		host:           config.Host,
		token:          config.Token,
		projects:       make(map[string][]*gitup.Repo),
		filterArchived: dd.Val(config.FilterArchived),
	}
	return g
}

func (g *gitlab) Projects() []*gitup.Repo {
	if len(g.projects) == 0 {
		g.fetchProjects()
	}
	result := []*gitup.Repo{}
	for _, v := range g.projects {
		result = append(result, v...)
	}
	return result
}

func (g *gitlab) ProjectsByGroup(group *string) ([]*gitup.Repo, error) {
	if len(g.projects) == 0 {
		g.fetchProjects()
	}
	// ). check if need to search subgroup
	prefix := dd.Val(group)
	subSearch := false
	if strings.Contains(prefix, "/") {
		prefix = prefix[:strings.IndexByte(prefix, '/')]
		subSearch = true
	}
	// ). find repos about target root group
	result, ok := g.projects[prefix]
	if !ok {
		return nil, fmt.Errorf("[GitLab]Not find projects in %s", *group)
	}
	if subSearch {
		// ). filter subgroup
		subResult := []*gitup.Repo{}
		for _, r := range result {
			if strings.HasPrefix(r.FullPath, dd.Val(group)) {
				subResult = append(subResult, r)
			}
		}
		if len(subResult) == 0 {
			return nil, fmt.Errorf("[GitLab]Not find projects in %s", *group)
		}
		result = subResult
	}
	return result, nil
}

const (
	perPage = 100
	baseURL = "https://%s/api/v4"
)

func (g *gitlab) fetchProjects() error {
	// ). init client
	gl, err := gitlabapi.NewClient(
		dd.Val(g.token),
		gitlabapi.WithBaseURL(fmt.Sprintf(baseURL, dd.Val(g.host))),
	)
	if err != nil {
		return err
	}

	// ). init context & channel
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	dst := make(chan []*gitlabapi.Project, 1)
	defer close(dst)

	go func() {
		defer cancel()

		// Prepare list project options
		opt := &gitlabapi.ListProjectsOptions{
			ListOptions: gitlabapi.ListOptions{
				Page:    1,
				PerPage: perPage,
			},
			Simple: dd.Ptr(true),
		}
		if g.filterArchived {
			opt.Archived = dd.Ptr(false)
		}

		for {
			// Get the first page with projects.
			ps, resp, err := gl.Projects.ListProjects(opt)
			if err != nil {
				infra.GetLogger().Log("[Gitlab]", "List projects error", err)
				return
			}

			dst <- ps

			// Exit the loop when we've seen all pages.
			if resp.NextPage == 0 {
				break
			}

			// Update the page number to get the next page.
			opt.Page = resp.NextPage
		}
	}()

	infra.GetLogger().Log("[Gitlab]", "Waiting fetching repo...")

	for alive := true; alive; {
		select {
		case ps := <-dst:
			convertToRepo(&g.projects, ps)

		case <-ctx.Done():
			infra.GetLogger().Log("[Gitlab]", "Done...")
			alive = false
		}
	}

	return nil
}

func convertToRepo(base *map[string][]*gitup.Repo, projects []*gitlabapi.Project) {
	for _, p := range projects {
		g := p.PathWithNamespace[:strings.IndexByte(p.PathWithNamespace, '/')]
		r := &gitup.Repo{
			URL:      p.HTTPURLToRepo,
			Name:     p.Name,
			Group:    g,
			FullPath: p.PathWithNamespace,
		}
		// fmt.Printf("%s - %s\n", r.Group, r.URL)
		ps, ok := (*base)[r.Group]
		if !ok {
			// the first repo insert about this group
			(*base)[r.Group] = append([]*gitup.Repo{}, r)
		} else {
			(*base)[r.Group] = append(ps, r)
		}
	}
}
