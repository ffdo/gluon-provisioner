package main

import (
	"github.com/BurntSushi/toml"
)

type Domain struct {
	Match  string
	Ignore bool
}

type Config struct {
	Domains map[string]Domain
}

func NewConfig(filename string) (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(filename, config)
	return
}
