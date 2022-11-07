package gitup

// Repo represent a repository
type Repo struct {
	ID       int
	URL      string
	Name     string
	Group    string
	FullPath string
}

// RepoListor represent a listor of all repositories
type RepoListor interface {
	Projects() []*Repo

	ProjectsByGroup(group *string) ([]*Repo, error)

	Project(group, name *string) (*Repo, error)
}

// RepoForker represent a forker to fork any repositories
type RepoForker interface {
	RepoListor

	Fork(r *Repo, group *string) (*Repo, error)

	Rename(r *Repo, name *string) (*Repo, error)

	Transfer(r *Repo, group *string) (*Repo, error)
}
