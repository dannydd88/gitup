package gitup

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"gitup/internal/infra"
	"gitup/pkg/git"

	"github.com/dannydd88/dd-go"
)

// Syncer -
type Syncer struct {
	Hub         RepoHub
	SyncConfig  *infra.SyncConfig
	Cwd         *string
	Concurrency int
	Logger      dd.Logger
}

// Go -
func (r *Syncer) Go() {
	r.Logger.Log("[Syncer]", "Started...")
	// ). Prepare repos
	var repos []*Repo
	if len(r.SyncConfig.Groups) == 0 {
		repos = r.Hub.Projects()
	} else {
		repos = []*Repo{}
		for _, g := range r.SyncConfig.Groups {
			result, err := r.Hub.ProjectsByGroup(g)
			if err != nil {
				r.Logger.Log("[Syncer]", "Meet error ->", err)
				continue
			} else {
				repos = append(repos, result...)
			}
		}
	}

	// ). Prepare context
	ctx, cancel := context.WithCancel(context.Background())
	output := make(chan string)
	defer close(output)

	// ). Prepare git repo
	var wg sync.WaitGroup

	r.Logger.Log("[Syncer]", "Start sync repos ->", len(repos))

	// ). Post git task to runner
	for _, repo := range repos {
		wg.Add(1)
		url := dd.String(repo.URL)
		path := dd.String(filepath.Join(dd.Val(r.Cwd), repo.FullPath))
		git := git.NewGit(r.Logger, url, path, r.SyncConfig.Bare)
		f := dd.Bind3(doSyncGitRepo, git, output, &wg)
		infra.GetWorkerPoolRunner().Post(f)
	}

	// ). Async wait task done
	go func() {
		defer cancel()
		r.Logger.Log("[Syncer]", "Waiting syncing repo...")
		wg.Wait()
	}()

	// ). Logging & wait all job done
	for alive := true; alive; {
		select {
		case m := <-output:
			r.Logger.Log(m)
		case <-ctx.Done():
			r.Logger.Log("[Syncer]", "Done...")
			alive = false
		}
	}
}

func doSyncGitRepo(g *git.Git, output chan string, wg *sync.WaitGroup) error {
	err := g.Sync()

	var msg string
	if err == nil {
		msg = fmt.Sprintf("[Syncer] Finish sync[%s]", dd.Val(g.Path()))
	} else {
		msg = fmt.Sprintf("[Syncer] Error sync[%s] err[%s]", dd.Val(g.Path()), err)
	}

	output <- msg
	wg.Done()

	return nil
}
