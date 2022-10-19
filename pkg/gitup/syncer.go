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

// Syncer
type Syncer struct {
	Api        RepoListor
	SyncConfig *infra.SyncConfig
	Cwd        *string
	Logger     dd.Logger
}

// Go
// Entrance of |Syncer|
func (s *Syncer) Go() {
	s.Logger.Log("[Syncer]", "Started...")

	// ). prepare repos
	var repos []*Repo
	if len(s.SyncConfig.Groups) == 0 {
		repos = s.Api.Projects()
	} else {
		repos = []*Repo{}
		for _, g := range s.SyncConfig.Groups {
			result, err := s.Api.ProjectsByGroup(g)
			if err != nil {
				s.Logger.Log("[Syncer]", "Meet error ->", err)
				continue
			} else {
				repos = append(repos, result...)
			}
		}
	}

	// ). prepare context
	ctx, cancel := context.WithCancel(context.Background())
	output := make(chan string)
	defer close(output)
	var wg sync.WaitGroup

	s.Logger.Log("[Syncer]", "Start sync repos ->", len(repos))

	// ). post git task to runner
	for _, repo := range repos {
		wg.Add(1)
		url := dd.String(repo.URL)
		path := dd.String(filepath.Join(dd.Val(s.Cwd), repo.FullPath))
		git := git.NewGit(s.Logger, url, path, s.SyncConfig.Bare)
		c := dd.Bind3(doSyncGitRepo, git, output, &wg)
		infra.GetWorkerPoolRunner().Post(c)
	}

	// ). async wait task done
	go func() {
		defer cancel()
		s.Logger.Log("[Syncer]", "Waiting syncing repo...")
		wg.Wait()
	}()

	// ). logging & wait all task done
	for alive := true; alive; {
		select {
		case m := <-output:
			s.Logger.Log(m)
		case <-ctx.Done():
			s.Logger.Log("[Syncer]", "Done...")
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

	return err
}
