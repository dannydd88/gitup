package gitup

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/dannydd88/gitup/pkg/git"

	"github.com/dannydd88/dd-go"
)

const (
	TagSync = "[sync]"
)

type SyncConfig struct {
	Token  *string
	Bare   bool
	Groups []*string
}

// Sync
type Sync struct {
	Api        RepoList
	SyncConfig *SyncConfig
	Cwd        *string
	TaskRunner dd.TaskRunner
	Logger     dd.LevelLogger
}

// Go
// Entrance of |sync|
func (s *Sync) Go() {
	s.Logger.Info(TagSync, "Started...")

	// ). prepare repos
	var repos []*Repo
	if len(s.SyncConfig.Groups) == 0 {
		repos = s.Api.Projects()
	} else {
		repos = []*Repo{}
		for _, g := range s.SyncConfig.Groups {
			result, err := s.Api.ProjectsByGroup(g)
			if err != nil {
				s.Logger.Warn(TagSync, "Meet error ->", err)
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
	wg := new(sync.WaitGroup)

	s.Logger.Info(TagSync, "Start sync repos ->", len(repos))

	// ). post git task to runner
	for _, repo := range repos {
		wg.Add(1)
		url := dd.Ptr(repo.URL)
		path := dd.Ptr(filepath.Join(dd.Val(s.Cwd), repo.FullPath))
		git := git.NewGoGit(s.Logger, &git.GitConfig{
			URL:     url,
			WorkDir: path,
			Bare:    s.SyncConfig.Bare,
			Token:   s.SyncConfig.Token,
		})
		c := dd.Bind3(doSyncGitRepo, git, output, wg)
		s.TaskRunner.Post(c)
	}

	// ). async wait task done
	go func() {
		defer cancel()
		s.Logger.Info(TagSync, "Waiting syncing repo...")
		wg.Wait()
	}()

	// ). logging & wait all task done
	for alive := true; alive; {
		select {
		case m := <-output:
			s.Logger.Info(m)
		case <-ctx.Done():
			s.Logger.Info(TagSync, "Done...")
			alive = false
		}
	}
}

func doSyncGitRepo(g git.Git, output chan string, wg *sync.WaitGroup) error {
	updated, err := g.Sync()

	var msg string
	if err == nil {
		var updateMsg string
		if updated {
			updateMsg = "sync-to-latest"
		} else {
			updateMsg = "already-up-to-date"
		}
		msg = fmt.Sprintf(
			"%s Finish sync[%s] [%s]",
			TagSync,
			dd.Val(g.Path()),
			updateMsg,
		)
	} else {
		msg = fmt.Sprintf(
			"%s Error sync[%s] err[%s]",
			TagSync,
			dd.Val(g.Path()),
			err,
		)
	}

	output <- msg
	wg.Done()

	return err
}
