package git

import (
	"io"
	"os"
	"path/filepath"

	"github.com/dannydd88/dd-go"
	gg "github.com/go-git/go-git/v5"
	gghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GoGit - a set of git commands via local cli git
type GoGit struct {
	config *GitConfig
	logger dd.Logger
}

// NewCLIGit - Init a new Git instance via cli git
func NewGoGit(logger dd.Logger, config *GitConfig) Git {
	// make sure |path| is exist
	if !dd.DirExists(config.WorkDir) {
		os.MkdirAll(dd.Val(config.WorkDir), os.ModePerm)
	}
	g := &GoGit{
		config: config,
		logger: logger,
	}
	return g
}

// Path - current git repo path
func (g *GoGit) Path() *string {
	return g.config.WorkDir
}

// Sync - Sync a git repository, clone if is a new one, update otherwise
func (g *GoGit) Sync() (bool, error) {
	var checkPath string
	path := dd.Val(g.config.WorkDir)
	if g.config.Bare {
		checkPath = filepath.Join(path, "HEAD")
	} else {
		checkPath = filepath.Join(path, ".git", "HEAD")
	}

	// update if repository already existed
	if dd.FileExists(dd.Ptr(checkPath)) {
		if g.config.Bare {
			return g.fetch()
		} else {
			return g.pull()
		}
	}
	// else clone
	return g.clone()
}

func (g *GoGit) clone() (bool, error) {
	path := dd.Val(g.config.WorkDir)
	g.logger.Log("[go-git]", "Clone repo ->", path)

	_, err := gg.PlainClone(path, g.config.Bare, &gg.CloneOptions{
		URL:      dd.Val(g.config.URL),
		Progress: io.Discard,
		Auth: &gghttp.BasicAuth{
			Username: "dummy",
			Password: dd.Val(g.config.Token),
		},
	})

	return err == nil, err
}

func (g *GoGit) fetch() (bool, error) {
	path := dd.Val(g.config.WorkDir)
	g.logger.Log("[go-git]", "fetch repo ->", path)

	r, err := gg.PlainOpen(path)
	if err != nil {
		return false, err
	}

	err = r.Fetch(&gg.FetchOptions{
		Progress: io.Discard,
		Auth: &gghttp.BasicAuth{
			Username: "dummy",
			Password: dd.Val(g.config.Token),
		},
	})
	return err == nil, err
}

func (g *GoGit) pull() (bool, error) {
	path := dd.Val(g.config.WorkDir)
	g.logger.Log("[go-git]", "pull repo ->", path)

	r, err := gg.PlainOpen(path)
	if err != nil {
		return false, err
	}

	w, err := r.Worktree()
	if err != nil {
		return false, err
	}

	err = w.Pull(&gg.PullOptions{
		Progress: io.Discard,
		Auth: &gghttp.BasicAuth{
			Username: "dummy",
			Password: dd.Val(g.config.Token),
		},
	})

	if err == gg.NoErrAlreadyUpToDate {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
