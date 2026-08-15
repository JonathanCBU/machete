package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	confPath      string = ".config/machete/conf.toml"
	defaultAwsDir string = ".aws"
	defaultSshDir string = ".ssh"
)

type Config struct {
	Aws string
	Ssh string
}

func GetConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	p := filepath.Join(home, confPath)

	// safe default handling if config file doesn't exist
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		log.Print("Loading default paths")
		return &Config{
			Aws: filepath.Join(home, defaultAwsDir),
			Ssh: filepath.Join(home, defaultSshDir),
		}, nil
	} else if err != nil {
		return nil, err
	}

	// read actual config file
	var confObj *Config
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	_, err = toml.Decode(string(data), &confObj)
	if err != nil {
		return nil, err
	}

	return confObj, nil
}
