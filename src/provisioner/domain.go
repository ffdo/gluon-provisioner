package main

import (
	"regexp"
)

func getDomain(nodeName string) (domain string, ignore bool, err error) {
	config, err := NewConfig(configFile)
	if err != nil {
		return
	}

	for domainName, domainConfig := range config.Domains {
		ignore = domainConfig.Ignore

		if domainConfig.Match != "" {
			var domainRe *regexp.Regexp
			domainRe, err = regexp.Compile(domainConfig.Match)
			if err != nil {
				return
			}

			if domainRe.MatchString(nodeName) {
				domain = domainName
				return
			}
		}

	}

	return
}
