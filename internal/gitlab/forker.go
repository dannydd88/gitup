package gitlab

import (
	"github.com/dannydd88/dd-go"
	"github.com/dannydd88/gitup/internal/infra"
	"github.com/dannydd88/gitup/pkg/gitup"
	gitlabapi "github.com/xanzy/go-gitlab"
)

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
