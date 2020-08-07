package config

import (
	"github.com/spf13/viper"

	"github.com/saltbo/authcar/pkg/oauth2"
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

type Config struct {
	Root    string        `yaml:"root"`
	Secret  string        `yaml:"secret"`
	Oauth2  oauth2.Config `yaml:"oauth2"`
	Roles   Roles         `yaml:"roles"`
	Routers []Router      `yaml:"routers"`
}

func Parse() (*Config, error) {
	conf := new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
