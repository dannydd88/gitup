package git

// GitConfig - configs relative with git
type GitConfig struct {
	URL     *string
	WorkDir *string
	Bare    bool
	Token   *string
}

// Git - a set of git commands to one git repository and one local path
type Git interface {
	// Path - current git repo path
	Path() *string

	// Sync - Sync a git repo, clone if is a new one, update otherwise
	//       |bool| indicate that whether repo is updated
	Sync() (bool, error)
}
