package config

import (
	"github.com/spf13/viper"
)

type Roles struct {
	Loader string `yaml:"loader"`
	//Upstream client.UpstreamConfig `yaml:"upstream"`
}

type Upstream struct {
	Address string            `yaml:"address"`
	Headers map[string]string `yaml:"headers"`
}

type Router struct {
	Pattern  string   `yaml:"pattern"`
	Upstream Upstream `yaml:"upstream"`
}

type Email struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Database struct {
	Driver string `yml:"driver"`
	DSN    string `yml:"dsn"`
}

type Routers []Router

type Config struct {
	SiteOrigin string   `yaml:"site_origin"`
	Database   Database `yml:"database"`
	Root       string   `yaml:"root"`
	Secret     string   `yaml:"secret"`
	Roles      Roles    `yaml:"roles"`
	Routers    Routers  `yaml:"routers"`
}

func Parse() (*Config, error) {
	conf := new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
