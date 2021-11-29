package gitup

// Repo represent a repository
type Repo struct {
	URL      string
	Name     string
	Group    string
	FullPath string
}

// RepoHub represent a hub of all repositories
type RepoHub interface {
	Projects() []*Repo

	ProjectsByGroup(group *string) ([]*Repo, error)
}
