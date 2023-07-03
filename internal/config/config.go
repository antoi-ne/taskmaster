package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Program stores the configuration attributes of an individual program in a config file.
type Program struct {
	Cmd          string            `yaml:"cmd"`
	NumProcs     uint              `yaml:"numprocs"`
	UMask        uint32            `yaml:"umask"`
	WorkingDir   string            `yaml:"workingdir"`
	AutoStart    bool              `yaml:"autostart"`
	AutoRestart  string            `yaml:"autorestart"`
	ExitCodes    []int             `yaml:"exitcodes"`
	StartRetries uint              `yaml:"startretries"`
	StartTime    uint              `yaml:"starttime"`
	StopSignal   string            `yaml:"stopsignal"`
	StopTime     uint              `yaml:"stoptime"`
	Stdout       string            `yaml:"stdout"`
	Stderr       string            `yaml:"stderr"`
	Env          map[string]string `yaml:"env"`
}

// Conf stores the configuration of multiple programs from a config file.
type Conf struct {
	Programs map[string]Program `yaml:"programs"`
}

// Parse reads the content of file name and tries to parse its content as yaml intoa config file structure.
func Parse(name string) (Conf, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return Conf{}, err
	}

	conf := Conf{}

	if err := yaml.Unmarshal(b, &conf); err != nil {
		return Conf{}, err
	}

	return conf, nil
}
