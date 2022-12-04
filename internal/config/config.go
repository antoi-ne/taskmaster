package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Parse reads the content of file name and tries to parse its content as yaml intoa config file structure
func Parse(name string) (*File, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	cf := new(File)

	if err := yaml.Unmarshal(b, cf); err != nil {
		return nil, err
	}

	return cf, nil
}
