package gitlab

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dannydd88/gitup/internal/infra"
	"github.com/dannydd88/gitup/pkg/gitup"

	"github.com/dannydd88/dd-go"
	gitlabapi "github.com/xanzy/go-gitlab"
)

const (
	perPage = 100
)

type gitlabListor struct {
	gitlabContext
	projects       map[string][]*gitup.Repo
	filterArchived bool
}

// NewListor
// Helper function to create a |RepoListor|'s gitlab implement
func NewListor(config *infra.RepoConfig) (gitup.RepoListor, error) {
	// ). construct gitlab client
	c, err := newGitlabClient(config.Token, config.Host)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabListor{
		gitlabContext: gitlabContext{
			apiClient: c,
		},
		projects:       make(map[string][]*gitup.Repo),
		filterArchived: config.FilterArchived,
	}
	return g, nil
}

func (g *gitlabListor) Projects() []*gitup.Repo {
	if len(g.projects) == 0 {
		g.fetchProjects()
	}
	result := []*gitup.Repo{}
	for _, v := range g.projects {
		result = append(result, v...)
	}
	return result
}

func (g *gitlabListor) ProjectsByGroup(group *string) ([]*gitup.Repo, error) {
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
		return nil, fmt.Errorf("[GitLab]Not find projects in %s", dd.Val(group))
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
			return nil, fmt.Errorf("[GitLab]Not find projects in %s", dd.Val(group))
		}
		result = subResult
	}
	return result, nil
}

func (g *gitlabListor) Project(group, name *string) (*gitup.Repo, error) {
	// ). list projects with group
	repos, err := g.ProjectsByGroup(group)
	if err != nil {
		return nil, err
	}

	for _, r := range repos {
		realGroup := r.FullPath[:strings.LastIndexByte(r.FullPath, '/')]
		if strings.Compare(dd.Val(group), realGroup) == 0 && strings.Compare(dd.Val(name), r.Name) == 0 {
			return r, nil
		}
	}

	return nil, fmt.Errorf("[Gitlab]Not find project[%s][%s]", dd.Val(group), dd.Val(name))
}

func (g *gitlabListor) fetchProjects() error {
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
			ps, resp, err := g.apiClient.Projects.ListProjects(opt)
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
			ID:       p.ID,
			URL:      p.HTTPURLToRepo,
			Name:     strings.TrimSpace(p.Name),
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

type gitlabForker struct {
	gitlabListor
	token string
	host  string
}

// NewForker
// Helper function to create |RepoForker|'s gitlab implement
func NewForker(config *infra.RepoConfig) (gitup.RepoForker, error) {
	// ). construct gitlab client
	c, err := newGitlabClient(config.Token, config.Host)
	if err != nil {
		return nil, err
	}

	// ). construct
	g := &gitlabForker{
		gitlabListor: gitlabListor{
			gitlabContext: gitlabContext{
				apiClient: c,
			},
			projects:       make(map[string][]*gitup.Repo),
			filterArchived: config.FilterArchived,
		},
		token: dd.Val(config.Token),
		host:  dd.Val(config.Host),
	}

	return g, nil
}

func (g *gitlabForker) Fork(r *gitup.Repo, group *string) (*gitup.Repo, error) {
	// ). prepare fork options
	opt := &gitlabapi.ForkProjectOptions{
		NamespacePath: group,
	}

	// ). do fork
	p, resp, err := g.apiClient.Projects.ForkProject(r.ID, opt)
	if err != nil {
		return nil, err
	}
	infra.GetLogger().Log("[Gitlab]", "Fork finish",
		"http ->", resp.StatusCode,
		"new project ->", p.ID,
	)

	// ). do disable project job token access
	{
		opt := &gitlabapi.PatchProjectJobTokenAccessSettingsOptions{
			Enabled: false,
		}
		resp, err := g.apiClient.JobTokenScope.PatchProjectJobTokenAccessSettings(p.ID, opt)
		if err != nil {
			return nil, err
		}
		infra.GetLogger().Log("[Gitlab]", "Disable project job token access",
			"http ->", resp.StatusCode,
		)
	}

	return &gitup.Repo{
		ID:       p.ID,
		Name:     p.Name,
		Group:    p.Namespace.FullPath,
		URL:      p.HTTPURLToRepo,
		FullPath: p.PathWithNamespace,
	}, nil
}

func (g *gitlabForker) Rename(r *gitup.Repo, name *string) (*gitup.Repo, error) {
	// ). prepare edit project options
	opt := &gitlabapi.EditProjectOptions{
		Name: name,
		Path: name,
	}

	// ). do rename
	p, resp, err := g.apiClient.Projects.EditProject(r.ID, opt)
	if err != nil {
		return nil, err
	}
	infra.GetLogger().Log("[Gitlab]", "Rename finish",
		"http ->", resp.StatusCode,
		"project ->", r.ID,
		"after ->", p.ID,
	)

	return &gitup.Repo{
		ID:       p.ID,
		Name:     p.Name,
		Group:    p.Namespace.FullPath,
		URL:      p.HTTPURLToRepo,
		FullPath: p.PathWithNamespace,
	}, nil
}

func (g *gitlabForker) Transfer(r *gitup.Repo, group *string) (*gitup.Repo, error) {
	// ). prepare transfer options
	opt := &gitlabapi.TransferProjectOptions{
		Namespace: group,
	}

	// ). do transfer
	p, resp, err := g.apiClient.Projects.TransferProject(r.ID, opt)
	if err != nil {
		return nil, err
	}
	infra.GetLogger().Log("[Gitlab]", "Transfer finish",
		"http ->", resp.StatusCode,
		"project ->", r.ID,
		"after ->", p.ID,
	)

	return &gitup.Repo{
		ID:       p.ID,
		Name:     p.Name,
		Group:    p.Namespace.FullPath,
		URL:      p.HTTPURLToRepo,
		FullPath: p.PathWithNamespace,
	}, nil
}

func (g *gitlabForker) DeleteForkRelationship(r *gitup.Repo) (bool, error) {
	// ). do delete fork relationship
	resp, err := g.apiClient.Projects.DeleteProjectForkRelation(r.ID)
	if err != nil {
		return false, err
	}
	infra.GetLogger().Log("[Gitlab]", "Delete fork relationship finish",
		"http -> ", resp.StatusCode,
		"project -> ", r.ID,
	)

	return true, nil
}
