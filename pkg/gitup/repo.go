package gitup

// Repo represent a repository
type Repo struct {
	ID       int
	URL      string
	Name     string
	Group    string
	FullPath string
}

// RepoList - represent a set of list operations of all repositories
type RepoList interface {
	Projects() []*Repo

	// ProjectsByGroup - List project by group name prefix match
	ProjectsByGroup(group *string) ([]*Repo, error)

	// Project - Filter target project by specific group and name
	Project(group, name *string) (*Repo, error)
}

// RepoFork - represent a set of fork operations to fork any repositories
type RepoFork interface {
	RepoList

	Fork(r *Repo, group *string) (*Repo, error)

	Rename(r *Repo, name *string) (*Repo, error)

	Transfer(r *Repo, group *string) (*Repo, error)

	DeleteForkRelationship(r *Repo) (bool, error)
}
