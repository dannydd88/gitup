package gitup

import (
	gitlabapi "github.com/xanzy/go-gitlab"
)

type gitlabFork struct {
	gitlabList
}

func (g *gitlabFork) Fork(r *Repo, group *string) (*Repo, error) {
	// ). prepare fork options
	opt := &gitlabapi.ForkProjectOptions{
		NamespacePath: group,
	}

	// ). do fork
	p, resp, err := g.Api().Projects.ForkProject(r.ID, opt)
	if err != nil {
		return nil, err
	}
	g.Logger().Info(
		TagGitlab,
		"Fork finish,",
		"http ->", resp.StatusCode,
		",",
		"new project ->", p.ID,
	)

	// ). do disable project job token access
	{
		opt := &gitlabapi.PatchProjectJobTokenAccessSettingsOptions{
			Enabled: false,
		}
		resp, err := g.Api().JobTokenScope.PatchProjectJobTokenAccessSettings(p.ID, opt)
		if err != nil {
			return nil, err
		}
		g.Logger().Info(
			TagGitlab,
			"Disable project job token access,",
			"http ->", resp.StatusCode,
		)
	}

	return &Repo{
		ID:       p.ID,
		Name:     p.Name,
		Group:    p.Namespace.FullPath,
		URL:      p.HTTPURLToRepo,
		FullPath: p.PathWithNamespace,
	}, nil
}

func (g *gitlabFork) Rename(r *Repo, name *string) (*Repo, error) {
	// ). prepare edit project options
	opt := &gitlabapi.EditProjectOptions{
		Name: name,
		Path: name,
	}

	// ). do rename
	p, resp, err := g.Api().Projects.EditProject(r.ID, opt)
	if err != nil {
		return nil, err
	}
	g.Logger().Info(
		TagGitlab,
		"Rename finish,",
		"http ->", resp.StatusCode,
		",",
		"project ->", r.ID,
		",",
		"after ->", p.ID,
	)

	return &Repo{
		ID:       p.ID,
		Name:     p.Name,
		Group:    p.Namespace.FullPath,
		URL:      p.HTTPURLToRepo,
		FullPath: p.PathWithNamespace,
	}, nil
}

func (g *gitlabFork) Transfer(r *Repo, group *string) (*Repo, error) {
	// ). prepare transfer options
	opt := &gitlabapi.TransferProjectOptions{
		Namespace: group,
	}

	// ). do transfer
	p, resp, err := g.Api().Projects.TransferProject(r.ID, opt)
	if err != nil {
		return nil, err
	}
	g.Logger().Info(
		TagGitlab,
		"Transfer finish,",
		"http ->", resp.StatusCode,
		",",
		"project ->", r.ID,
		",",
		"after ->", p.ID,
	)

	return &Repo{
		ID:       p.ID,
		Name:     p.Name,
		Group:    p.Namespace.FullPath,
		URL:      p.HTTPURLToRepo,
		FullPath: p.PathWithNamespace,
	}, nil
}

func (g *gitlabFork) DeleteForkRelationship(r *Repo) (bool, error) {
	// ). do delete fork relationship
	resp, err := g.Api().Projects.DeleteProjectForkRelation(r.ID)
	if err != nil {
		return false, err
	}
	g.Logger().Info(
		TagGitlab,
		"Delete fork relationship finish,",
		"http -> ", resp.StatusCode,
		",",
		"project -> ", r.ID,
	)

	return true, nil
}
