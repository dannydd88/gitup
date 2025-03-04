package gitup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dannydd88/dd-go"
	gitlabapi "gitlab.com/gitlab-org/api/client-go"
)

const (
	perPage = 100
)

type gitlabList struct {
	GitlabApi
	filterArchived bool
}

func (g *gitlabList) Projects() []*Repo {
	// ). prepare result
	result := []*Repo{}

	// ). fetch all projects
	projects, err := g.fetchProjects(nil)
	if err != nil {
		return result
	}

	// ). build final result
	for _, v := range dd.Val(projects) {
		result = append(result, v...)
	}
	return result
}

func (g *gitlabList) ProjectsByGroup(group *string) ([]*Repo, error) {
	// ). fetch projects by |group| string
	projects, err := g.fetchProjects(group)
	if err != nil {
		return nil, err
	}

	// ). check if need to search subgroup
	prefix := dd.Val(group)
	subSearch := false
	if strings.Contains(prefix, "/") {
		prefix = prefix[:strings.IndexByte(prefix, '/')]
		subSearch = true
	}

	// ). find repos about target root group
	result, ok := dd.Val(projects)[prefix]
	if !ok {
		return nil, fmt.Errorf("[gitlab] Not find projects in %s", dd.Val(group))
	}
	if subSearch {
		// ). filter subgroup
		subResult := []*Repo{}
		for _, r := range result {
			if strings.HasPrefix(r.FullPath, dd.Val(group)) {
				subResult = append(subResult, r)
			}
		}
		if len(subResult) == 0 {
			return nil, fmt.Errorf("[gitlab] Not find projects in %s", dd.Val(group))
		}
		result = subResult
	}
	return result, nil
}

func (g *gitlabList) Project(group, name *string) (*Repo, error) {
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

	return nil, fmt.Errorf("[gitlab] Not find project[%s][%s]", dd.Val(group), dd.Val(name))
}

func (g *gitlabList) fetchProjects(group *string) (*map[string][]*Repo, error) {
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
		if group != nil {
			opt.SearchNamespaces = dd.Ptr(true)
			opt.Search = group
		}

		for {
			// Get the first page with projects.
			ps, resp, err := g.Api().Projects.ListProjects(opt)
			if err != nil {
				g.Logger().Info(TagGitlab, "List projects error", err)
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

	g.Logger().Info(
		TagGitlab,
		"Waiting fetching repo...",
		"[", "Group", "->", dd.Val(group), "]",
	)
	output := make(map[string][]*Repo)

	for alive := true; alive; {
		select {
		case ps := <-dst:
			convertToRepo(dd.Ptr(output), ps)

		case <-ctx.Done():
			g.Logger().Info(TagGitlab, "Done...")
			alive = false
		}
	}

	return dd.Ptr(output), nil
}

func convertToRepo(base *map[string][]*Repo, projects []*gitlabapi.Project) {
	for _, p := range projects {
		g := p.PathWithNamespace[:strings.IndexByte(p.PathWithNamespace, '/')]
		r := &Repo{
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
			(*base)[r.Group] = append([]*Repo{}, r)
		} else {
			(*base)[r.Group] = append(ps, r)
		}
	}
}
