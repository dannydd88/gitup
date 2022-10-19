package gitup

import (
	"context"
	"fmt"
	"sync"

	"github.com/dannydd88/dd-go"
)

type ForkConfig struct {
	FromGroup *string   `yaml:"from-group"`
	FromRepos []*string `yaml:"from-repos"`
	ToGroup   *string   `yaml:"to-group"`
	ToRepos   []*string `yaml:"to-repos,omitempty"`
}

// Forker
type Forker struct {
	Api         RepoForker
	ForkConfigs []*ForkConfig
	TaskRunner  dd.TaskRunner
	Logger      dd.Logger
}

type forkDetail struct {
	source      *Repo
	targetGroup *string
	targetName  *string
}

// Go
// Entrance of |Forker|
func (f *Forker) Go() {
	f.Logger.Log("[Forker]", "Started...")

	// ). prepare context
	ctx, cancel := context.WithCancel(context.Background())
	output := make(chan string)
	defer close(output)
	var wg sync.WaitGroup

	// ). do fork in each |ForkConfig|
	for _, fc := range f.ForkConfigs {
		// ). check config
		if len(fc.ToRepos) != 0 && len(fc.FromRepos) != len(fc.ToRepos) {
			f.Logger.Log("[Forker]", "find to-repos != from-repos error in from-group ->", fc.FromGroup, ", skip this!")
			continue
		}

		// ). foreach target repo
		for i, r := range fc.FromRepos {
			// ). find target repo
			repo, err := f.Api.Project(fc.FromGroup, r)
			if err != nil {
				f.Logger.Log("[Forker]", "finding source repo meet error ->", err)
				continue
			}

			// ). prepare fork detail
			detail := &forkDetail{
				source:      repo,
				targetGroup: fc.ToGroup,
			}
			if len(fc.ToRepos) != 0 {
				detail.targetName = fc.ToRepos[i]
			}
			wg.Add(1)

			// ). async do fork
			c := dd.Bind4(doFork, f.Api, detail, output, &wg)
			f.TaskRunner.Post(c)
		}
	}

	// ). async wait task done
	go func() {
		defer cancel()
		f.Logger.Log("[Forker]", "Waiting forking repo...")
		wg.Wait()
	}()

	// ). logging & wait all task done
	for alive := true; alive; {
		select {
		case m := <-output:
			f.Logger.Log(m)
		case <-ctx.Done():
			f.Logger.Log("[Forker]", "Done...")
			alive = false
		}
	}
}

func doFork(api RepoForker, detail *forkDetail, output chan string, wg *sync.WaitGroup) error {
	err := api.Fork(detail.source, detail.targetGroup, detail.targetName)

	var msg string
	if err == nil {
		msg = fmt.Sprintf(
			"[Forker] Fork success [%s]->[%s][%s]",
			detail.source.FullPath,
			dd.Val(detail.targetGroup),
			dd.Val(detail.targetName),
		)
	} else {
		msg = fmt.Sprintf(
			"[Forker] Fork error [%s]->[%s][%s] err[%s]",
			detail.source.FullPath,
			dd.Val(detail.targetGroup),
			dd.Val(detail.targetName),
			err,
		)
	}

	output <- msg
	wg.Done()

	return err
}
