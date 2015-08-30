package main

import (
	"regexp"
	"strings"
)

func getDomain(nodeName string) (domain string, ignore bool, err error) {
	config, err := NewConfig(configFile)
	if err != nil {
		return
	}

	for domainName, domainConfig := range config.Domains {
		ignore = domainConfig.Ignore

		for _, re := range domainConfig.Matches {
			if re != "" {
				var domainRe *regexp.Regexp
				domainRe, err = regexp.Compile(strings.ToLower(re))
				if err != nil {
					return
				}

				if domainRe.MatchString(strings.ToLower(nodeName)) {
					domain = domainName
					return
				}
			}
		}

	}

	return
}
