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
	Hub         RepoHub
	Git         GitConfig
	Cwd         string
	Concurrency int
	Logger      base.Logger
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

	// ). Prepare context
	ctx, cancel := context.WithCancel(context.Background())

	// ). Prepare git and ready to send
	input := make(chan *Git)
	defer close(input)
	var wg sync.WaitGroup
	go func() {
		r.Logger.Log("[Runner]Start sync repos ->", len(repos))
		for _, repo := range repos {
			wg.Add(1)
			url := base.String(repo.URL)
			path := base.String(filepath.Join(r.Cwd, repo.FullPath))
			input <- NewGit(r.Logger, url, path, base.Bool(r.Git.Bare))
		}
	}()

	<-time.After(1 * time.Second)

	// ). Prepare taskRunners
	output := make(chan string)
	defer close(output)
	taskRunners := make([]*taskRunner, r.Concurrency)
	for i := 0; i < r.Concurrency; i++ {
		taskRunners[i] = &taskRunner{
			input:  input,
			output: output,
			wg:     &wg,
			ctx:    ctx,
		}
		taskRunners[i].run()
	}

	go func() {
		defer cancel()
		r.Logger.Log("[Runner]Waiting syncing repo...")
		wg.Wait()
	}()

	stop := false
	for !stop {
		select {
		case m := <-output:
			r.Logger.Log(m)
		case <-ctx.Done():
			r.Logger.Log("[Runner]Done...")
			stop = true
			break
		}
	}
}

type taskRunner struct {
	input  chan *Git
	output chan string
	wg     *sync.WaitGroup
	ctx    context.Context
}

func (t *taskRunner) run() {
	go func() {
		stop := false
		for !stop {
			select {
			case g := <-t.input:
				err := g.Sync()
				t.output <- fmt.Sprintf("[Runner]Finish sync[%s] with error[%s]",
					base.StringValue(g.path), err)
				t.wg.Done()
			case <-t.ctx.Done():
				stop = true
			}
		}
	}()
}
