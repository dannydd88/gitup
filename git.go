package gitup

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codeskyblue/go-sh"
	"github.com/dannydd88/gitup/base"
)

// Git represent a set of git commands to one git repository and one local path
type Git struct {
	sess   *sh.Session
	url    *string
	path   *string
	bare   *bool
	logger base.Logger
}

// NewGit - Init a new Git instance
func NewGit(logger base.Logger, url, path *string, bare *bool) *Git {
	// make sure |path| is exist
	if !base.DirExists(path) {
		os.MkdirAll(*path, os.ModePerm)
	}
	g := &Git{
		sess:   sh.NewSession(),
		url:    url,
		path:   path,
		bare:   bare,
		logger: logger,
	}
	g.sess.Stdout = ioutil.Discard
	g.sess.Stderr = ioutil.Discard
	g.sess.SetDir(*g.path)
	return g
}

// Sync - Sync a git repository, clone if is a new one, update otherwise
func (g *Git) Sync() error {
	var checkPath string
	if base.BoolValue(g.bare) {
		checkPath = filepath.Join(base.StringValue(g.path), "HEAD")
	} else {
		checkPath = filepath.Join(base.StringValue(g.path), ".git", "HEAD")
	}

	// update if repository already existed
	if base.FileExists(base.String(checkPath)) {
		return g.Update()
	}
	// else clone
	return g.Clone()
}

// Clone - clone a new git repository
func (g *Git) Clone() error {
	g.logger.Log("[Git]Clone repo ->", base.StringValue(g.path))
	params := []string{"clone"}
	if *g.bare {
		params = append(params, "--bare")
	}
	params = append(params, *g.url, *g.path)
	return g.sess.Command("git", params).Run()
}

// Update - update a git repository
func (g *Git) Update() error {
	g.logger.Log("[Git]Update repo ->", base.StringValue(g.path))
	var p string
	if base.BoolValue(g.bare) {
		p = "fetch"
	} else {
		p = "pull"
	}
	return g.sess.Command("git", p).Run()
}
