package config

import (
	"fmt"
	"log"

	"github.com/saltbo/gopkg/mailutil"
	"github.com/spf13/viper"

	"github.com/saltbo/moreu/pkg/gormutil"
)

type (
	Router struct {
		Pattern  string `yaml:"pattern" json:"pattern"`
		Upstream struct {
			Address string            `yaml:"address" json:"address"`
			Headers map[string]string `yaml:"headers" json:"headers"`
		} `yaml:"upstream" json:"upstream"`
	}
	Static struct {
		Pattern string `yaml:"pattern" json:"pattern"`
		DistDir string `yaml:"distdir" json:"dist_dir"`
	}
)

type Config struct {
	Debug      bool            `yaml:"debug" json:"debug"`
	Secret     string          `yaml:"secret" json:"secret"`
	MoreuRoot  string          `yaml:"moreu_root" json:"moreu_root"`
	GRbacFile  string          `yaml:"grbac_file" json:"grbac_file"`
	Invitation bool            `yaml:"invitation" json:"invitation"`
	Email      mailutil.Config `yaml:"email" json:"email"`
	Database   gormutil.Config `yaml:"database" json:"database"`
	Statics    []Static        `yaml:"statics" json:"statics"`
	Routers    []Router        `yaml:"routers" json:"routers"`

	OuterKeys []string `yaml:"outer_keys" json:"outer_keys"`
}

func (c *Config) EmailAct() bool {
	return c.Email.Host != ""
}

func Parse() *Config {
	conf := new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		log.Fatalln(err)
	}

	fmt.Print(conf)
	return conf
}
