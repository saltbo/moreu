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
	Static struct {
		Pattern string `yaml:"pattern"`
		DistDir string `yaml:"distdir"`
	}
)

type Config struct {
	Debug      bool            `yaml:"debug"`
	Secret     string          `yaml:"secret"`
	MoreuRoot  string          `yaml:"moreu_root"`
	Invitation bool            `yaml:"invitation"`
	Email      mailutil.Config `yaml:"email"`
	Database   gormutil.Config `yaml:"database"`
	Statics    []Static        `yaml:"statics"`
	Routers    []Router        `yaml:"routers"`
}

func Parse() (*Config, error) {
	conf := new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
