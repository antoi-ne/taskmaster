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

	f := new(File)

	if err := yaml.Unmarshal(b, f); err != nil {
		return nil, err
	}

	return f, nil
}
