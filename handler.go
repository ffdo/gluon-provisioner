package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func PathHandler(path string, pathConfig *PathConfig, nodeDB *NodeDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		remoteIP := net.ParseIP(req.Header.Get(*xffHeader))
		if remoteIP == nil {
			msg := fmt.Sprint("%s: Cannot parse IP address in %s header", req.URL.Path, *xffHeader)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		var node *Node
	RulesLoop:
		for _, rule := range pathConfig.Rules {
			if node == nil && (len(rule.When) != 0 || rule.Careful) {
				node = nodeDB.GetNode(remoteIP)
				if node == nil {
					log.Printf("%s %s: Node not found in node DB, cannot evaluate rule", remoteIP, req.URL.Path)
					continue
				}
				log.Printf("%s %s: Identified node %s", remoteIP, req.URL.Path, node.Nodeinfo.Hostname)
			}

			for _, condition := range rule.When {
				if !condition.Check(node.Nodeinfo) {
					continue RulesLoop
				}
				log.Printf("%s %s: Field '%s' matches '%s'", remoteIP, req.URL.Path, condition.Field, condition.Match)
			}

			if rule.Careful && !node.HasOnlyVPNLinks() {
				log.Printf("%s %s: Careful enabled and node is not VPN only, skipping rule", remoteIP, req.URL.Path)
				continue
			}

			if rule.Disabled {
				log.Printf("%s %s: Matching rule is DISABLED: %s -> %s", remoteIP, req.URL.Path, path, rule.Path)
				continue
			}

			log.Printf("%s %s: Using redirect rule %s -> %s", remoteIP, req.URL.Path, path, rule.Path)
			w.Header().Set(*redirectHeader, fmt.Sprint(rule.Path, strings.TrimPrefix(req.URL.Path, path)))
			return
		}

		log.Printf("%s %s: Using default redirect %s -> %s", remoteIP, req.URL.Path, path, pathConfig.Default)
		w.Header().Set(*redirectHeader, fmt.Sprint(pathConfig.Default, strings.TrimPrefix(req.URL.Path, path)))
		return
	})
}
