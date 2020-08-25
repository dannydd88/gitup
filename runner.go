package gitup

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/dannydd88/gitup/base"
)

// Runner -
type Runner struct {
	Hub    RepoHub
	Git    GitConfig
	Cwd    string
	Logger base.Logger
}

// Execute -
func (r *Runner) Execute() {
	// ). Prepare repos
	var repos []*Repo
	if len(r.Git.Groups) == 0 {
		repos = r.Hub.Projects()
	} else {
		repos = []*Repo{}
		for _, g := range r.Git.Groups {
			result, err := r.Hub.ProjectsByGroup(base.String(g))
			if err != nil {
				continue
			} else {
				repos = append(repos, result...)
			}
		}
	}

	// ). Sync repos
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	dst := make(chan string)
	defer close(dst)
	{
		count := len(repos)
		var wg sync.WaitGroup
		wg.Add(count)
		r.Logger.Log("[Runner]Start sync repos ->", count)
		for _, repo := range repos {
			url := base.String(repo.URL)
			path := base.String(filepath.Join(r.Cwd, repo.FullPath))
			g := NewGit(r.Logger, url, path, base.Bool(r.Git.Bare))
			go func(g *Git) {
				defer wg.Done()

				err := g.Sync()
				dst <- fmt.Sprintf("[Runner]Finish sync[%s] with error[%s]",
					base.StringValue(g.path), err)
			}(g)
		}

		go func(ctx context.Context) {
			defer cancel()

			r.Logger.Log("[Runner]Waiting syncing repo...")
			wg.Wait()
		}(ctx)
	}

	stop := false
	for !stop {
		select {
		case m := <-dst:
			r.Logger.Log(m)
		case <-ctx.Done():
			r.Logger.Log("[Runner]Done...")
			stop = true
			break
		}
	}
}
