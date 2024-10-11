package infra

// RepoConfig - repo setion of config.yaml
type RepoConfig struct {
	Type           *string `yaml:"type"`
	Host           *string `yaml:"host,omitempty"`
	Token          *string `yaml:"token,omitempty"`
	FilterArchived bool    `yaml:"filter_archived,omitempty"`
}

// SyncConfig - sync setion of config.yaml
type SyncConfig struct {
	Bare   bool      `yaml:"bare"`
	Groups []*string `yaml:"groups,omitempty"`
}

// Config - config represent config.yaml
type Config struct {
	RepoConfig *RepoConfig `yaml:"repo"`
	SyncConfig *SyncConfig `yaml:"sync"`
	Cwd        *string     `yaml:"cwd"`
}

// INIConfig - config represent gitup.ini
type INIConfig struct {
	Host  *string `ini:"gitlab_host"`
	Token *string `ini:"gitlab_token"`
}
