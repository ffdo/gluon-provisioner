package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	networks       map[string]*Network
	sortedNetworks []*Network
}

func NewConfig(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	networks := make(map[string]*Network)
	err = yaml.Unmarshal(content, networks)
	if err != nil {
		return nil, err
	}

	c := &Config{
		networks:       networks,
		sortedNetworks: make([]*Network, 0, len(networks)),
	}

	for networkString, networkConfig := range networks {
		_, networkConfig.ipnet, err = net.ParseCIDR(networkString)
		if err != nil {
			return nil, err
		}
		networkConfig.serveMux = http.NewServeMux()
		networkConfig.nodeDB = NewNodeDB(*updateInterval, networkConfig.Nodes, networkConfig.Graph)

		for path, pathConfig := range networkConfig.Routes {
			for _, rule := range pathConfig.Rules {
				for _, condition := range rule.When {
					condition.re, err = regexp.Compile(strings.ToLower(condition.Match))
					if err != nil {
						return nil, err
					}
				}
			}
			networkConfig.serveMux.Handle(path, PathHandler(path, pathConfig, networkConfig.nodeDB))
		}

		c.sortedNetworks = append(c.sortedNetworks, networkConfig)
	}

	sort.Sort(ByNetmask(c.sortedNetworks))

	return c, nil
}

type Network struct {
	Nodes, Graph string
	Routes       map[string]*PathConfig

	ipnet    *net.IPNet     `yaml:"-"`
	nodeDB   *NodeDB        `yaml:"-"`
	serveMux *http.ServeMux `yaml:"-"`
}

type PathConfig struct {
	Default string
	Rules   []*Rule
}

type Rule struct {
	When     []*Condition
	Path     string
	Careful  bool
	Disabled bool
}

type ByNetmask []*Network

// Sort ByNetmask largest to smallest
func (a ByNetmask) Len() int           { return len(a) }
func (a ByNetmask) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNetmask) Less(i, j int) bool { return bytes.Compare(a[i].ipnet.IP, a[j].ipnet.IP) > 0 }

func (c *Config) GetServeMux(ip net.IP) *http.ServeMux {
	for _, network := range c.sortedNetworks {
		if network.ipnet.Contains(ip) {
			return network.serveMux
		}
	}
	return nil
}
