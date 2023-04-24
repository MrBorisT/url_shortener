package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	MaxShortUrlLen int    `yaml:"max_short_url_len"`
	Port           string `yaml:"port"`
	Redis          struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		TTL      int    `yaml:"ttl_in_minutes"`
	} `yaml:"redis"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return errors.WithMessage(err, "parsing yaml")
	}

	return nil
}
