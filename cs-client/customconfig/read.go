package customconfig

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Read(path string) (*Config, error) {
	var config Config
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
