package main

import (
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

type Domain struct {
	Matches []string
	Ignore  bool
	Force   []string
}

type Config struct {
	Domains map[string]Domain
}

func NewConfig(filename string) (c *Config, err error) {
	c = &Config{}
	_, err = toml.DecodeFile(filename, c)
	return
}

func (c *Config) getDomain(node *Node) (domain string, force bool, ignore bool, err error) {
	for domainName, domainConfig := range c.Domains {
		ignore = domainConfig.Ignore

		// Check if node should be forced
		for _, nodespec := range domainConfig.Force {
			if nodespec != "" {
				nodespec = strings.ToLower(nodespec)
				if nodespec == strings.ToLower(node.Nodeinfo.Hostname) {
					domain = domainName
					force = true
					return
				}
				for _, ip := range node.Nodeinfo.Network.Addresses {
					if nodespec == ip {
						domain = domainName
						force = true
						return
					}
				}
			}
		}

		for _, re := range domainConfig.Matches {
			if re != "" {
				var domainRe *regexp.Regexp
				domainRe, err = regexp.Compile(strings.ToLower(re))
				if err != nil {
					return
				}

				if domainRe.MatchString(strings.ToLower(node.Nodeinfo.Hostname)) {
					domain = domainName
					return
				}
			}
		}

	}

	return
}
