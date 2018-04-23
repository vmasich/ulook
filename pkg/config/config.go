package config

import (
	"github.com/BurntSushi/toml"
	"github.com/coreos/pkg/capnslog"
)

type Config struct {
	NatsURL string
}

var conf Config

var log = capnslog.NewPackageLogger("bitbucket.org/vmasych/urllookup/pkg/config", "config")

func init() {
	New()
}

func New() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
}

func Get() *Config {
	return &conf
}
