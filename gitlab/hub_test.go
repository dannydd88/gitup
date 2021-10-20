package gitlab

import (
	"testing"

	"github.com/dannydd88/gitup"
)

func TestConstruct(t *testing.T) {
	c := &gitup.RepoConfig{
		Host:  "aaa.com",
		Token: "bbb",
	}
	r := NewGitlab(c)
	nr, ok := r.(*gitlab)
	if !ok || nr == nil {
		t.Errorf("construct error %v", r)
	}
}

const projectJson = `
{
	"id": 1094,
	"description": "",
	"name": "work-ticket",
	"name_with_namespace": "a-b / b-online / work-ticket",
	"path": "work-ticket",
	"path_with_namespace": "a-b/b-online/work-ticket",
	"created_at": "2020-05-07T16:37:54.490+08:00",
	"default_branch": "master",
	"tag_list": [],
	"ssh_url_to_repo": "git@git.xxx.com:a-b/b-online/work-ticket.git",
	"http_url_to_repo": "https://git.xxx.com/a-b/b-online/work-ticket.git",
	"web_url": "https://git.xxx.com/a-b/b-online/work-ticket",
	"readme_url": "https://git.xxx.com/a-b/b-online/work-ticket/blob/master/README.md",
	"avatar_url": null,
	"star_count": 0,
	"forks_count": 0,
	"last_activity_at": "2020-07-02T17:08:36.988+08:00",
	"namespace": {
			"id": 290,
			"name": "b-online",
			"path": "b-online",
			"kind": "group",
			"full_path": "a-b/b-online",
			"parent_id": 269,
			"avatar_url": null,
			"web_url": "https://git.xxx.com/groups/a-b/b-online"
	},
	"packages_enabled": true,
	"empty_repo": false,
	"archived": false,
	"visibility": "private"
}`
