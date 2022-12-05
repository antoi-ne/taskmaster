package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Program stores the configuration attributes of an individual program in a config file.
type Program struct {
	Cmd          string            `yaml:"cmd"`
	NumProcs     uint              `yaml:"nuprocs"`
	UMask        uint32            `yaml:"umask"`
	WorkingDir   string            `yaml:"workingdir"`
	AutoStart    bool              `yaml:"autostart"`
	AutoRestart  string            `yaml:"autorestart"`
	ExitCode     []int             `yaml:"exitcodes"`
	StartRetries uint              `yaml:"startretries"`
	StartTime    uint              `yaml:"starttime"`
	StopSignal   string            `yaml:"stopsignal"`
	StopTime     uint              `yaml:"stoptime"`
	Stdout       string            `yaml:"stdout"`
	Stderr       string            `yaml:"stderr"`
	Env          map[string]string `yaml:"env"`
}

// File stores the configuration of multiple programs from a config file.
type File struct {
	Programs map[string]*Program `yaml:"programs"`
}

// Parse reads the content of file name and tries to parse its content as yaml intoa config file structure.
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
