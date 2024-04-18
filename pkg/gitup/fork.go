package gitup

import (
	"context"
	"fmt"
	"sync"

	"github.com/dannydd88/dd-go"
)

type ForkConfig struct {
	FromGroup      *string   `yaml:"from-group"`
	FromRepos      []*string `yaml:"from-repos"`
	ToGroup        *string   `yaml:"to-group"`
	ToRepos        []*string `yaml:"to-repos,omitempty"`
	RmForkRelation *bool     `yaml:"rm-fork-relation"`
}

// Fork
type Fork struct {
	Api         RepoFork
	ForkConfigs []*ForkConfig
	TaskRunner  dd.TaskRunner
	Logger      dd.LevelLogger
}

type forkDetail struct {
	source         *Repo
	targetGroup    *string
	targetName     *string
	sameGroupFork  bool
	changeNameFork bool
	rmForkRelation bool
}

// Go
// Entrance of |fork|
func (f *Fork) Go() {
	f.Logger.Log("[fork]", "Started...")

	// ). prepare context
	ctx, cancel := context.WithCancel(context.Background())
	output := make(chan string)
	defer close(output)
	wg := new(sync.WaitGroup)

	// ). do fork in each |ForkConfig|
	for _, fc := range f.ForkConfigs {
		// ). check config
		if len(fc.ToRepos) != 0 && len(fc.FromRepos) != len(fc.ToRepos) {
			f.Logger.Warn(
				"[fork]",
				"find len(to-repos) != len(from-repos) error in from-group ->",
				fc.FromGroup,
				", skip this!",
			)
			continue
		}

		// ). foreach target repo
		for i, r := range fc.FromRepos {
			// ). find target repo
			repo, err := f.Api.Project(fc.FromGroup, r)
			if err != nil {
				f.Logger.Warn("[fork]", "finding source repo meet error ->", err)
				continue
			}

			// ). prepare fork detail
			detail := &forkDetail{
				source:         repo,
				rmForkRelation: dd.Val(fc.RmForkRelation),
			}
			if fc.ToGroup == nil {
				detail.targetGroup = fc.FromGroup
			} else {
				detail.targetGroup = fc.ToGroup
			}
			if len(fc.ToRepos) != 0 {
				detail.targetName = fc.ToRepos[i]
			}
			if dd.Val(fc.FromGroup) == dd.Val(detail.targetGroup) {
				detail.sameGroupFork = true
				if detail.targetName == nil {
					f.Logger.Warn(
						"[fork]",
						"same group fork [",
						detail.source.Name,
						"] without new repo name, skip this",
					)
					continue
				}
			}
			if detail.targetName != nil && dd.Val(detail.targetName) != dd.Val(r) {
				detail.changeNameFork = true
			}
			wg.Add(1)

			// ). async do fork
			c := dd.Bind4(doFork, f.Api, detail, output, wg)
			f.TaskRunner.Post(c)
		}
	}

	// ). async wait task done
	go func() {
		defer cancel()
		f.Logger.Log("[fork]", "Waiting forking repo...")
		wg.Wait()
	}()

	// ). logging & wait all task done
	for alive := true; alive; {
		select {
		case m := <-output:
			f.Logger.Log(m)
		case <-ctx.Done():
			f.Logger.Log("[fork]", "Done...")
			alive = false
		}
	}
}

func doFork(api RepoFork, detail *forkDetail, output chan string, wg *sync.WaitGroup) error {
	// ). do fork
	targetGroup := detail.targetGroup
	if detail.sameGroupFork {
		targetGroup = nil
	}
	forkedRepo, err := api.Fork(detail.source, targetGroup)

	// ). do rename if necessary
	if err == nil && detail.changeNameFork {
		_, err = api.Rename(forkedRepo, detail.targetName)
	}

	// ). do transfer if necessary
	if err == nil && detail.sameGroupFork {
		_, err = api.Transfer(forkedRepo, detail.targetGroup)
	}

	// ). do remove fork relationship if necessary
	if err == nil && detail.rmForkRelation {
		_, err = api.DeleteForkRelationship(forkedRepo)
	}

	var msg string
	if err == nil {
		msg = fmt.Sprintf(
			"[fork] Fork success [%s]->[%s][%s]",
			detail.source.FullPath,
			dd.Val(detail.targetGroup),
			dd.Val(detail.targetName),
		)
	} else {
		msg = fmt.Sprintf(
			"[fork] Fork error [%s]->[%s][%s] err[%s]",
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
