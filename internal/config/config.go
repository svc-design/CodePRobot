package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	WatchPaths      []string `yaml:"watch_paths"`
	TriggerKeywords []string `yaml:"trigger_keywords"`
	OpenAI          struct {
		Model       string  `yaml:"model"`
		Temperature float64 `yaml:"temperature"`
	} `yaml:"openai"`
	GitHub struct {
		Repo     string   `yaml:"repo"`
		Token    string   `yaml:"token"`
		Reviewer []string `yaml:"reviewer"`
	} `yaml:"github"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

