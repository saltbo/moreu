package config

import (
	"github.com/spf13/viper"

	"github.com/saltbo/moreu/pkg/mailutil"
	"github.com/saltbo/moreu/pkg/ormutil"
)

type (
	Router struct {
		Pattern  string `yaml:"pattern"`
		Upstream struct {
			Address string            `yaml:"address"`
			Headers map[string]string `yaml:"headers"`
		} `yaml:"upstream"`
	}
	Routers []Router
)

type Config struct {
	Host     string          `yaml:"host"`
	Root     string          `yaml:"root"`
	Secret   string          `yaml:"secret"`
	Email    mailutil.Config `yaml:"email"`
	Database ormutil.Config  `yaml:"database"`
	Routers  Routers         `yaml:"routers"`
}

func Parse() (*Config, error) {
	conf := new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
