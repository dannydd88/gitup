package git

import (
	"io"
	"os"
	"path/filepath"

	"github.com/codeskyblue/go-sh"
	"github.com/dannydd88/dd-go"
)

// CLIGit - a set of git commands via local cli git
type CLIGit struct {
	sess   *sh.Session
	config *GitConfig
	logger dd.Logger
}

// NewCLIGit - Init a new Git instance via cli git
func NewCLIGit(logger dd.Logger, config *GitConfig) Git {
	// make sure |path| is exist
	if !dd.DirExists(config.WorkDir) {
		os.MkdirAll(dd.Val(config.WorkDir), os.ModePerm)
	}
	g := &CLIGit{
		sess:   sh.NewSession(),
		config: config,
		logger: logger,
	}
	g.sess.Stdout = io.Discard
	g.sess.Stderr = io.Discard
	g.sess.SetDir(dd.Val(g.config.WorkDir))
	return g
}

// Path - current git repo path
func (g *CLIGit) Path() *string {
	return g.config.WorkDir
}

// Sync - Sync a git repository, clone if is a new one, update otherwise
func (g *CLIGit) Sync() error {
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

func (g *CLIGit) clone() error {
	path := dd.Val(g.config.WorkDir)
	g.logger.Log("[cli-git]", "Clone repo ->", path)
	params := []string{"clone"}
	if g.config.Bare {
		params = append(params, "--bare")
	}
	params = append(params, dd.Val(g.config.URL), path)
	return g.sess.Command("git", params).Run()
}

func (g *CLIGit) fetch() error {
	g.logger.Log("[cli-git]", "fetch repo ->", dd.Val(g.config.WorkDir))
	return g.sess.Command("git", "fetch").Run()
}

func (g *CLIGit) pull() error {
	g.logger.Log("[cli-git]", "pull repo ->", dd.Val(g.config.WorkDir))
	return g.sess.Command("git", "pull").Run()
}
