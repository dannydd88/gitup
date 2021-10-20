package gitup

import (
	"encoding/json"
	"io/ioutil"
)

// RepoConfig - repo setion of config.json
type RepoConfig struct {
	Type           string `json:"type"`
	Host           string `json:"host"`
	Token          string `json:"token"`
	FilterArchived bool   `json:"filter_archived,omitempty"`
}

// GitConfig - git setion of config.json
type GitConfig struct {
	Bare   bool     `json:"bare"`
	Groups []string `json:"groups,omitempty"`
}

// Config - config represent config.json
type Config struct {
	RepoConfig RepoConfig `json:"repo"`
	GitConfig  GitConfig  `json:"git"`
	Cwd        string     `json:"cwd"`
}

// LoadConfig - Load config
func LoadConfig(filepath *string) (*Config, error) {
	c := Config{}
	data, err := ioutil.ReadFile(*filepath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
