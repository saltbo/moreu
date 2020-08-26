package config

import (
	"github.com/saltbo/gopkg/gormutil"
	"github.com/saltbo/gopkg/mailutil"
	"github.com/spf13/viper"
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

type (
	Static struct {
		Pattern string `yaml:"pattern"`
		DistDir string `yaml:"distdir"`
	}
	Statics []Static
)

type Config struct {
	Host     string          `yaml:"host"`
	Moreu    string          `yaml:"moreu"`
	Secret   string          `yaml:"secret"`
	Invite   bool            `yaml:"invite"`
	Email    mailutil.Config `yaml:"email"`
	Database gormutil.Config `yaml:"database"`
	Statics  Statics         `yaml:"statics"`
	Routers  Routers         `yaml:"routers"`
}

func Parse() (*Config, error) {
	conf := new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
