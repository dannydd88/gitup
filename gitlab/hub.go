package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dannydd88/gitup"
	"github.com/dannydd88/gitup/base"
)

type gitlab struct {
	host     *string
	token    *string
	projects map[string][]*gitup.Repo
}

// NewGitlab -
func NewGitlab(host, token *string) gitup.RepoHub {
	g := &gitlab{
		host:     host,
		token:    token,
		projects: make(map[string][]*gitup.Repo),
	}
	return g
}

func (g *gitlab) Projects() []*gitup.Repo {
	if len(g.projects) == 0 {
		g.fetchProjects()
	}
	result := []*gitup.Repo{}
	for _, v := range g.projects {
		result = append(result, v...)
	}
	return result
}

func (g *gitlab) ProjectsByGroup(group *string) ([]*gitup.Repo, error) {
	if len(g.projects) == 0 {
		g.fetchProjects()
	}
	result, ok := g.projects[base.StringValue(group)]
	if !ok {
		return nil, fmt.Errorf("[GitLab]Not find projects in %s", *group)
	}
	return result, nil
}

const (
	perPage    = 100
	projectURL = "https://%s/api/v4/projects?private_token=%s&page=%d&per_page=%d&simple=true"
)

type project struct {
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	RepoHTTPURL       string `json:"http_url_to_repo"`
}

func (g *gitlab) fetchProjects() error {
	resp, err := httpRequest(g.host, g.token, 1, perPage)
	if err != nil {
		return err
	}

	totalPage, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
	if err != nil {
		return err
	}

	body, err := readResponse(resp)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	dst := make(chan []byte, 1)
	defer close(dst)
	{
		var wg sync.WaitGroup
		wg.Add(totalPage - 1)
		dst <- body

		for i := 2; i <= totalPage; i++ {
			go func(page int) {
				defer wg.Done()

				resp, err = httpRequest(g.host, g.token, page, perPage)
				if err != nil {
					return
				}
				body, err := readResponse(resp)
				if err != nil {
					return
				}

				// fmt.Printf("Fetched for page[%d]\n", page)
				dst <- body
			}(i)
		}

		go func() {
			defer cancel()

			fmt.Printf("[GitLab]Waiting fetching repo...\n")
			wg.Wait()
		}()
	}

	var lastError error
	stop := false
	for !stop {
		select {
		case b := <-dst:
			v := []project{}
			err = json.Unmarshal(b, &v)
			if err != nil {
				stop = true
				lastError = err
				break
			}
			convertToRepo(&g.projects, &v)
		case <-ctx.Done():
			fmt.Printf("[Gitlab]Done...\n")
			stop = true
		}
	}

	return lastError
}

func httpRequest(host, token *string, page, perPage int) (*http.Response, error) {
	url := fmt.Sprintf(projectURL,
		base.StringValue(host),
		base.StringValue(token),
		page,
		perPage)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func readResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func convertToRepo(base *map[string][]*gitup.Repo, projects *[]project) {
	for _, p := range *projects {
		g := p.PathWithNamespace[:strings.IndexByte(p.PathWithNamespace, '/')]
		r := &gitup.Repo{
			URL:      p.RepoHTTPURL,
			Name:     p.Name,
			Group:    g,
			FullPath: p.PathWithNamespace,
		}
		// fmt.Printf("%s - %s\n", r.Group, r.URL)
		ps, ok := (*base)[r.Group]
		if !ok {
			(*base)[r.Group] = append([]*gitup.Repo{}, r)
		} else {
			(*base)[r.Group] = append(ps, r)
		}
	}
}
