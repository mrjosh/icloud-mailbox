package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ConfMap struct {
	Environment string      `yaml:"environment"`
	Auth        AuthConfig  `yaml:"auth"`
	SMTP        *SMTPConfig `yaml:"smtp"`
}

type AuthConfig struct {
	SecretID string `yaml:"secret_id"`
}

type SMTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	From string `yaml:"from"`
	Name string `yaml:"name"`
}

func LoadFile(filename string) (*ConfMap, error) {
	appConfig := new(ConfMap)
	confOut, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	if err := yaml.NewDecoder(confOut).Decode(appConfig); err != nil {
		return nil, err
	}
	return appConfig, nil
}
