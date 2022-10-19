package infra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dannydd88/dd-go"
	"gopkg.in/yaml.v3"
)

// RepoConfig - repo setion of config.yaml
type RepoConfig struct {
	Type           *string `yaml:"type"`
	Host           *string `yaml:"host"`
	Token          *string `yaml:"token"`
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

// LoadConfig - Load config
func loadConfig(path *string) (*Config, error) {
	// ). find out final full config filepath
	var p string
	if filepath.IsAbs(*path) {
		p = *path
	} else {
		dir, err := os.Getwd()
		if err == nil {
			p = filepath.Join(dir, *path)
		}
	}
	if !dd.FileExists(dd.Ptr(p)) {
		return nil, fmt.Errorf("cannot find config -> %s", p)
	}

	// ). do load config
	c := Config{}
	data, err := os.ReadFile(dd.Val(path))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
